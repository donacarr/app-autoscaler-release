package app_test

import (
	. "acceptance/helpers"
	"fmt"
	"os"

	"github.com/KevinJCross/cf-test-helpers/v2/cf"
	"github.com/KevinJCross/cf-test-helpers/v2/generator"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	. "github.com/onsi/gomega/gexec"

	"time"
)

var _ = Describe("AutoScaler recurring schedule policy", func() {
	var (
		appName              string
		appGUID              string
		initialInstanceCount int
		daysOfMonthOrWeek    Days
		startTime            time.Time
		endTime              time.Time
		policy               string
	)

	BeforeEach(func() {
		if cfg.IsServiceOfferingEnabled() {
			instanceName = generator.PrefixedRandomName(cfg.Prefix, cfg.InstancePrefix)
			createService := cf.Cf("create-service", cfg.ServiceName, cfg.ServicePlan, instanceName, "-b", cfg.ServiceBroker).Wait(cfg.DefaultTimeoutDuration())
			Expect(createService).To(Exit(0), "failed creating service")
		}

		initialInstanceCount = 1
		appName = CreateTestApp(cfg, "recurring-schedule", initialInstanceCount)
		appGUID = GetAppGuid(cfg, appName)
	})

	AfterEach(func() {
		if os.Getenv("SKIP_TEARDOWN") == "true" {
			fmt.Println("Skipping Teardown...")
		} else {
			DeletePolicy(appName, appGUID)
			Expect(cf.Cf("delete", appName, "-f", "-r").Wait(cfg.DefaultTimeoutDuration())).To(Exit(0))
		}
	})

	Context("when scaling by recurring schedule", func() {

		JustBeforeEach(func() {
			location, err := time.LoadLocation("GMT")
			Expect(err).NotTo(HaveOccurred())
			startTime, endTime = getStartAndEndTime(location, 70*time.Second, time.Duration(interval+120)*time.Second)
			scheduleInstanceInit := 3
			policy = GenerateDynamicAndRecurringSchedulePolicy(1, 4, 80, "GMT", startTime, endTime, daysOfMonthOrWeek, 2, 5, scheduleInstanceInit)
			instanceName = CreatePolicy(cfg, appName, appGUID, policy)
			Expect(cf.Cf("start", appName).Wait(cfg.CfPushTimeoutDuration())).To(Exit(0))
			WaitForNInstancesRunning(appGUID, scheduleInstanceInit, cfg.DefaultTimeoutDuration())
		})

		Context("with days of month", func() {
			BeforeEach(func() {
				daysOfMonthOrWeek = DaysOfMonth
			})

			It("should scale", func() {
				By("setting to initial_min_instance_count")
				jobRunTime := time.Until(startTime.Add(5 * time.Minute))
				WaitForNInstancesRunning(appGUID, 3, jobRunTime)

				By("setting schedule's instance_min_count")
				jobRunTime = time.Until(endTime)
				Eventually(func() int {
					return RunningInstances(appGUID, jobRunTime)
				}, jobRunTime, 15*time.Second).Should(Equal(2))

				jobRunTime = time.Until(endTime)
				Consistently(func() int {
					return RunningInstances(appGUID, jobRunTime)
				}, jobRunTime, 15*time.Second).Should(Equal(2))

				By("setting to default instance_min_count")
				WaitForNInstancesRunning(appGUID, 1, time.Duration(interval+60)*time.Second)
			})

		})

		Context("with days of week", func() {
			BeforeEach(func() {
				daysOfMonthOrWeek = DaysOfWeek
			})

			It("should scale", func() {
				By("setting to initial_min_instance_count")
				jobRunTime := time.Until(startTime.Add(5 * time.Minute))
				WaitForNInstancesRunning(appGUID, 3, jobRunTime)

				By("setting schedule's instance_min_count")
				jobRunTime = time.Until(endTime)
				Eventually(func() int {
					return RunningInstances(appGUID, jobRunTime)
				}, jobRunTime, 15*time.Second).Should(Equal(2))

				jobRunTime = time.Until(endTime)
				Consistently(func() int {
					return RunningInstances(appGUID, jobRunTime)
				}, jobRunTime, 15*time.Second).Should(Equal(2))

				By("setting to default instance_min_count")
				WaitForNInstancesRunning(appGUID, 1, time.Duration(interval+120)*time.Second)
			})
		})
	})

})
