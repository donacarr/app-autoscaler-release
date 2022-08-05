package cf

import (
	"fmt"
	"net/url"
)

const (
	RoleOrganisationUser           RoleType = "organization_user"
	RoleOrganizationAuditor        RoleType = "organization_auditor"
	RoleOrganizationManager        RoleType = "organization_manager"
	RoleOrganizationBillingManager RoleType = "organization_billing_manager"
	RoleSpaceAuditor               RoleType = "space_auditor"
	RoleSpaceDeveloper             RoleType = "space_developer"
	RoleSpaceManager               RoleType = "space_manager"
	RoleSpaceSupporter             RoleType = "space_supporter"
)

type (
	Role struct {
		Guid string   `json:"guid"`
		Type RoleType `json:"type"`
	}
	RoleType string
)

type Roles []Role

func (r Roles) HasRole(roleType RoleType) bool {
	for _, role := range r {
		if role.Type == roleType {
			return true
		}
	}
	return false
}

/*GetSpaceDeveloperRoles
 * Get role information given a set of filters
 * from the v3 api https://v3-apidocs.cloudfoundry.org/version/3.122.0/index.html#roles
 */
func (c *Client) GetSpaceDeveloperRoles(spaceId string, userId string) (Roles, error) {
	parameters := url.Values{}
	parameters.Add("types", "space_developer")
	parameters.Add("space_guids", spaceId)
	parameters.Add("user_guids", userId)
	params := parameters.Encode()
	theUrl := fmt.Sprintf("%s/v3/roles?%s", c.conf.API, params)

	return ResourceRetriever[Role]{c}.getAllPages(theUrl)
}
