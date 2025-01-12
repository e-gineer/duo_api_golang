package admin

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	duoapi "github.com/duosecurity/duo_api_golang"
)

// Client provides access to Duo's admin API.
type Client struct {
	duoapi.DuoApi
}

type ListResultMetadata struct {
	NextOffset   json.Number `json:"next_offset"`
	PrevOffset   json.Number `json:"prev_offset"`
	TotalObjects json.Number `json:"total_objects"`
}

type ListResult struct {
	Metadata ListResultMetadata `json:"metadata"`
}

func (l *ListResult) metadata() ListResultMetadata {
	return l.Metadata
}

// New initializes an admin API Client struct.
func New(base duoapi.DuoApi) *Client {
	return &Client{base}
}

// User models a single user.
type User struct {
	Alias1            *string           `json:"alias1" url:"alias1"`
	Alias2            *string           `json:"alias2" url:"alias2"`
	Alias3            *string           `json:"alias3" url:"alias3"`
	Alias4            *string           `json:"alias4" url:"alias4"`
	Aliases           map[string]string `json:"aliases" url:"aliases"`
	Created           uint64            `json:"created"`
	Email             string            `json:"email" url:"email"`
	FirstName         *string           `json:"firstname" url:"firstname"`
	Groups            []Group           `json:"groups"`
	IsEnrolled        bool              `json:"is_enrolled"`
	LastDirectorySync *uint64           `json:"last_directory_sync"`
	LastLogin         *uint64           `json:"last_login"`
	LastName          *string           `json:"lastname" url:"lastname"`
	Notes             string            `json:"notes" url:"notes"`
	Phones            []Phone           `json:"phones"`
	RealName          *string           `json:"realname" url:"realname"`
	Status            string            `json:"status" url:"status"`
	Tokens            []Token           `json:"tokens"`
	UserID            string            `json:"user_id"`
	Username          string            `json:"username" url:"username"`
}

// URLValues transforms a User into url.Values using the 'url' struct tag to
// define the key of the map. Fields are skiped if the value is empty.
func (u *User) URLValues() url.Values {
	params := url.Values{}

	t := reflect.TypeOf(u).Elem()
	v := reflect.ValueOf(u).Elem()

	// Iterate over all available struct fields
	for i := 0; i < t.NumField(); i++ {
		tag := t.Field(i).Tag.Get("url")
		if tag == "" {
			continue
		}
		// Skip fields have a zero value.
		if v.Field(i).Interface() == reflect.Zero(v.Field(i).Type()).Interface() {
			continue
		}
		var val string
		if t.Field(i).Type.Kind() == reflect.Ptr {
			val = fmt.Sprintf("%v", v.Field(i).Elem())
		} else {
			val = fmt.Sprintf("%v", v.Field(i))
		}
		params[tag] = []string{val}
	}
	return params
}

// Group models a group to which users may belong.
type Group struct {
	Desc             string `json:"desc"`
	GroupID          string `json:"group_id"`
	MobileOTPEnabled bool   `json:"mobile_otp_enabled"`
	Name             string `json:"name"`
	PushEnabled      bool   `json:"push_enabled"`
	SMSEnabled       bool   `json:"sms_enabled"`
	Status           string `json:"status"`
	VoiceEnabled     bool   `json:"voice_enabled"`
}

// Phone models a user's phone.
type Phone struct {
	Activated        bool     `json:"activated"`
	Capabilities     []string `json:"capabilities"`
	Encrypted        string   `json:"encrypted"`
	Extension        string   `json:"extension"`
	Fingerprint      string   `json:"fingerprint"`
	LastSeen         string   `json:"last_seen"`
	Model            string   `json:"model"`
	Name             string   `json:"name"`
	Number           string   `json:"number"`
	PhoneID          string   `json:"phone_id"`
	Platform         string   `json:"platform"`
	Postdelay        string   `json:"postdelay"`
	Predelay         string   `json:"predelay"`
	Screenlock       string   `json:"screenlock"`
	SMSPasscodesSent bool     `json:"sms_passcodes_sent"`
	Tampered         string   `json:"tampered"`
	Type             string   `json:"type"`
	Users            []User   `json:"users"`
}

// Token models a hardware security token.
type Token struct {
	TokenID  string `json:"token_id"`
	Type     string `json:"type"`
	Serial   string `json:"serial"`
	TOTPStep *int   `json:"totp_step"`
	Users    []User `json:"users"`
}

// U2FToken models a U2F security token.
type U2FToken struct {
	DateAdded      uint64 `json:"date_added"`
	RegistrationID string `json:"registration_id"`
	User           *User  `json:"user"`
}

// Application integrations (not including SSO, which is excluded by the API)
type Integration struct {
	AdminApiAdmins             int      `json:"adminapi_admins"`
	AdminApiInfo               int      `json:"adminapi_info"`
	AdminApiIntegrations       int      `json:"adminapi_integrations"`
	AdminApiReadLog            int      `json:"adminapi_read_log"`
	AdminApiReadResource       int      `json:"adminapi_read_resource"`
	AdminApiSettings           int      `json:"adminapi_settings"`
	AdminApiWriteResource      int      `json:"adminapi_write_resource"`
	FramelessAuthPromptEnabled *int     `json:"frameless_auth_prompt_enabled,omitempty"`
	Greeting                   string   `json:"greeting"`
	GroupsAllowed              []string `json:"groups_allowed"`
	IntegrationKey             string   `json:"integration_key"`
	Name                       string   `json:"name"`
	NetworksForApiAccess       *string  `json:"networks_for_api_access,omitempty"`
	Notes                      string   `json:"notes"`
	PolicyKey                  string   `json:"policy_key,omitempty"`
	PromptV4Enabled            string   `json:"prompt_v4_enabled"`
	SecretKey                  string   `json:"secret_key"`
	// Note: API says int of 1 or false, not sure how to handle here?
	SelfServiceAllowed          interface{} `json:"self_service_allowed,omitempty"`
	Type                        string      `json:"type"`
	UsernameNormalizationPolicy string      `json:"username_normalization_policy"`
}

// Administrator models an admin user.
type Administrator struct {
	AdminID                string   `json:"admin_id"`
	AdminUnits             []string `json:"admin_units"`
	Created                uint64   `json:"created"`
	Email                  string   `json:"email"`
	HardToken              *Token   `json:"hardtoken"`
	LastLogin              uint64   `json:"last_login"`
	Name                   string   `json:"name"`
	PasswordChangeRequired bool     `json:"password_change_required"`
	Phone                  string   `json:"phone"`
	RestrictedByAdminUnits bool     `json:"restricted_by_admin_units"`
	Role                   string   `json:"role"`
	Status                 string   `json:"status"`
}

// AdministrativeUnit models an administrative unit.
type AdministrativeUnit struct {
	AdminUnitID            string    `json:"admin_unit_id"`
	Description            string    `json:"description"`
	Name                   string    `json:"name"`
	Groups                 *[]string `json:"groups"`
	Integrations           *[]string `json:"integrations"`
	RestrictByGroups       bool      `json:"restrict_by_groups"`
	RestrictByIntegrations bool      `json:"restrict_by_integrations"`
}

// Common URL options

// Limit sets the optional limit parameter for an API request.
func Limit(limit uint64) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("limit", strconv.FormatUint(limit, 10))
	}
}

// Offset sets the optional offset parameter for an API request.
func Offset(offset uint64) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("offset", strconv.FormatUint(offset, 10))
	}
}

// User methods

// GetUsersUsername sets the optional username parameter for a GetUsers request.
func GetUsersUsername(name string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("username", name)
	}
}

// GetUsersResult models responses containing a list of users.
type GetUsersResult struct {
	duoapi.StatResult
	ListResult
	Response []User
}

// GetUserResult models responses containing a single user.
type GetUserResult struct {
	duoapi.StatResult
	Response User
}

func (result *GetUsersResult) getResponse() interface{} {
	return result.Response
}

func (result *GetUsersResult) appendResponse(users interface{}) {
	asserted_users := users.([]User)
	result.Response = append(result.Response, asserted_users...)
}

// GetUsers calls GET /admin/v1/users
// See https://duo.com/docs/adminapi#retrieve-users
func (c *Client) GetUsers(options ...func(*url.Values)) (*GetUsersResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	cb := func(params url.Values) (responsePage, error) {
		return c.retrieveUsers(params)
	}
	response, err := c.retrieveItems(params, cb)
	if err != nil {
		return nil, err
	}

	return response.(*GetUsersResult), nil
}

type responsePage interface {
	metadata() ListResultMetadata
	getResponse() interface{}
	appendResponse(interface{})
}

type pageFetcher func(params url.Values) (responsePage, error)

func (c *Client) retrieveItems(
	params url.Values,
	fetcher pageFetcher,
) (responsePage, error) {
	if params.Get("offset") == "" {
		params.Set("offset", "0")
	}

	if params.Get("limit") == "" {
		params.Set("limit", "100")
		accumulator, firstErr := fetcher(params)

		if firstErr != nil {
			return nil, firstErr
		}

		params.Set("offset", accumulator.metadata().NextOffset.String())
		for params.Get("offset") != "" {
			nextResult, err := fetcher(params)
			if err != nil {
				return nil, err
			}
			nextResult.appendResponse(accumulator.getResponse())
			accumulator = nextResult
			params.Set("offset", accumulator.metadata().NextOffset.String())
		}
		return accumulator, nil
	}

	return fetcher(params)
}

func (c *Client) retrieveUsers(params url.Values) (*GetUsersResult, error) {
	_, body, err := c.SignedCall(http.MethodGet, "/admin/v1/users", params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetUsersResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetUser calls GET /admin/v1/users/:user_id
// See https://duo.com/docs/adminapi#retrieve-user-by-id
func (c *Client) GetUser(userID string) (*GetUserResult, error) {
	path := fmt.Sprintf("/admin/v1/users/%s", userID)

	_, body, err := c.SignedCall(http.MethodGet, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetUserResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// CreateUser calls POST /admin/v1/users
// See https://duo.com/docs/adminapi#create-user
func (c *Client) CreateUser(params url.Values) (*GetUserResult, error) {
	path := "/admin/v1/users"

	_, body, err := c.SignedCall(http.MethodPost, path, params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetUserResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// ModifyUser calls POST /admin/v1/users/:user_id
// See https://duo.com/docs/adminapi#modify-user
func (c *Client) ModifyUser(userID string, params url.Values) (*GetUserResult, error) {
	path := fmt.Sprintf("/admin/v1/users/%s", userID)

	_, body, err := c.SignedCall(http.MethodPost, path, params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetUserResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeleteUser calls DELETE /admin/v1/users/:user_id
// See https://duo.com/docs/adminapi#delete-user
func (c *Client) DeleteUser(userID string) (*duoapi.StatResult, error) {
	path := fmt.Sprintf("/admin/v1/users/%s", userID)

	_, body, err := c.SignedCall(http.MethodDelete, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &duoapi.StatResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetUserGroups calls GET /admin/v1/users/:user_id/groups
// See https://duo.com/docs/adminapi#retrieve-groups-by-user-id
func (c *Client) GetUserGroups(userID string, options ...func(*url.Values)) (*GetGroupsResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	cb := func(params url.Values) (responsePage, error) {
		return c.retrieveUserGroups(userID, params)
	}
	response, err := c.retrieveItems(params, cb)
	if err != nil {
		return nil, err
	}

	return response.(*GetGroupsResult), nil
}

// AssociateGroupWithUser calls POST /admin/v1/users/:user_id/groups
// See https://duo.com/docs/adminapi#associate-group-with-user
func (c *Client) AssociateGroupWithUser(userID string, groupID string) (*duoapi.StatResult, error) {
	path := fmt.Sprintf("/admin/v1/users/%s/groups", userID)

	params := url.Values{}
	params.Set("group_id", groupID)

	_, body, err := c.SignedCall(http.MethodPost, path, params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &duoapi.StatResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DisassociateGroupFromUser calls POST /admin/v1/users/:user_id/groups
// See https://duo.com/docs/adminapi#disassociate-group-from-user
func (c *Client) DisassociateGroupFromUser(userID string, groupID string) (*duoapi.StatResult, error) {
	path := fmt.Sprintf("/admin/v1/users/%s/groups/%s", userID, groupID)

	_, body, err := c.SignedCall(http.MethodDelete, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &duoapi.StatResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (c *Client) retrieveUserGroups(userID string, params url.Values) (*GetGroupsResult, error) {
	path := fmt.Sprintf("/admin/v1/users/%s/groups", userID)

	_, body, err := c.SignedCall(http.MethodGet, path, params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetGroupsResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetUserPhones calls GET /admin/v1/users/:user_id/phones
// See https://duo.com/docs/adminapi#retrieve-phones-by-user-id
func (c *Client) GetUserPhones(userID string, options ...func(*url.Values)) (*GetPhonesResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	cb := func(params url.Values) (responsePage, error) {
		return c.retrieveUserPhones(userID, params)
	}
	response, err := c.retrieveItems(params, cb)
	if err != nil {
		return nil, err
	}

	return response.(*GetPhonesResult), nil
}

func (c *Client) retrieveUserPhones(userID string, params url.Values) (*GetPhonesResult, error) {
	path := fmt.Sprintf("/admin/v1/users/%s/phones", userID)

	_, body, err := c.SignedCall(http.MethodGet, path, params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetPhonesResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetUserTokens calls GET /admin/v1/users/:user_id/tokens
// See https://duo.com/docs/adminapi#retrieve-hardware-tokens-by-user-id
func (c *Client) GetUserTokens(userID string, options ...func(*url.Values)) (*GetTokensResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	cb := func(params url.Values) (responsePage, error) {
		return c.retrieveUserTokens(userID, params)
	}
	response, err := c.retrieveItems(params, cb)
	if err != nil {
		return nil, err
	}

	return response.(*GetTokensResult), nil
}

func (c *Client) retrieveUserTokens(userID string, params url.Values) (*GetTokensResult, error) {
	path := fmt.Sprintf("/admin/v1/users/%s/tokens", userID)

	_, body, err := c.SignedCall(http.MethodGet, path, params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetTokensResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// StringResult models responses containing a simple string.
type StringResult struct {
	duoapi.StatResult
	Response string
}

// AssociateUserToken calls POST /admin/v1/users/:user_id/tokens
// See https://duo.com/docs/adminapi#associate-hardware-token-with-user
func (c *Client) AssociateUserToken(userID, tokenID string) (*StringResult, error) {
	path := fmt.Sprintf("/admin/v1/users/%s/tokens", userID)

	params := url.Values{}
	params.Set("token_id", tokenID)

	_, body, err := c.SignedCall(http.MethodPost, path, params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &StringResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetUserU2FTokens calls GET /admin/v1/users/:user_id/u2ftokens
// See https://duo.com/docs/adminapi#retrieve-u2f-tokens-by-user-id
func (c *Client) GetUserU2FTokens(userID string, options ...func(*url.Values)) (*GetU2FTokensResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	cb := func(params url.Values) (responsePage, error) {
		return c.retrieveUserU2FTokens(userID, params)
	}
	response, err := c.retrieveItems(params, cb)
	if err != nil {
		return nil, err
	}

	return response.(*GetU2FTokensResult), nil
}

func (c *Client) retrieveUserU2FTokens(userID string, params url.Values) (*GetU2FTokensResult, error) {
	path := fmt.Sprintf("/admin/v1/users/%s/u2ftokens", userID)

	_, body, err := c.SignedCall(http.MethodGet, path, params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetU2FTokensResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// StringArrayResult models response containing an array of strings.
type StringArrayResult struct {
	duoapi.StatResult
	Response []string
}

// GetUserBypassCodes calls POST /admin/v1/users/:user_id/bypass_codes
// see https://duo.com/docs/adminapi#create-bypass-codes-for-user
func (c *Client) GetUserBypassCodes(userID string, options ...func(*url.Values)) (*StringArrayResult, error) {
	path := fmt.Sprintf("/admin/v1/users/%s/bypass_codes", userID)

	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	_, body, err := c.SignedCall(http.MethodPost, path, params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &StringArrayResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// Group methods

// GetGroupsResult models responses containing a list of groups.
type GetGroupsResult struct {
	duoapi.StatResult
	ListResult
	Response []Group
}

func (result *GetGroupsResult) getResponse() interface{} {
	return result.Response
}

func (result *GetGroupsResult) appendResponse(groups interface{}) {
	asserted_groups := groups.([]Group)
	result.Response = append(result.Response, asserted_groups...)
}

// GetGroups calls GET /admin/v1/groups
// See https://duo.com/docs/adminapi#retrieve-groups
func (c *Client) GetGroups(options ...func(*url.Values)) (*GetGroupsResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	cb := func(params url.Values) (responsePage, error) {
		return c.retrieveGroups(params)
	}
	response, err := c.retrieveItems(params, cb)
	if err != nil {
		return nil, err
	}

	return response.(*GetGroupsResult), nil
}

func (c *Client) retrieveGroups(params url.Values) (*GetGroupsResult, error) {
	_, body, err := c.SignedCall(http.MethodGet, "/admin/v1/groups", params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetGroupsResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetGroupResult models responses containing a single group.
type GetGroupResult struct {
	duoapi.StatResult
	Response Group
}

// GetGroup calls GET /admin/v2/group/:group_id
// See https://duo.com/docs/adminapi#get-group-info
func (c *Client) GetGroup(groupID string) (*GetGroupResult, error) {
	path := fmt.Sprintf("/admin/v2/groups/%s", groupID)

	_, body, err := c.SignedCall(http.MethodGet, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetGroupResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Phone methods

// GetPhonesNumber sets the optional number parameter for a GetPhones request.
func GetPhonesNumber(number string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("number", number)
	}
}

// GetPhonesExtension sets the optional extension parameter for a GetPhones request.
func GetPhonesExtension(ext string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("extension", ext)
	}
}

// GetPhonesResult models responses containing a list of phones.
type GetPhonesResult struct {
	duoapi.StatResult
	ListResult
	Response []Phone
}

func (result *GetPhonesResult) getResponse() interface{} {
	return result.Response
}

func (result *GetPhonesResult) appendResponse(phones interface{}) {
	asserted_phones := phones.([]Phone)
	result.Response = append(result.Response, asserted_phones...)
}

// GetPhones calls GET /admin/v1/phones
// See https://duo.com/docs/adminapi#phones
func (c *Client) GetPhones(options ...func(*url.Values)) (*GetPhonesResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	cb := func(params url.Values) (responsePage, error) {
		return c.retrievePhones(params)
	}
	response, err := c.retrieveItems(params, cb)
	if err != nil {
		return nil, err
	}

	return response.(*GetPhonesResult), nil
}

func (c *Client) retrievePhones(params url.Values) (*GetPhonesResult, error) {
	_, body, err := c.SignedCall(http.MethodGet, "/admin/v1/phones", params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetPhonesResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetPhoneResult models responses containing a single phone.
type GetPhoneResult struct {
	duoapi.StatResult
	Response Phone
}

// GetPhone calls GET /admin/v1/phones/:phone_id
// See https://duo.com/docs/adminapi#retrieve-phone-by-id
func (c *Client) GetPhone(phoneID string) (*GetPhoneResult, error) {
	path := fmt.Sprintf("/admin/v1/phones/%s", phoneID)

	_, body, err := c.SignedCall(http.MethodGet, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetPhoneResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// DeletePhone calls DELETE /admin/v1/phones/:phone_id
// See https://duo.com/docs/adminapi#delete-phone
func (c *Client) DeletePhone(phoneID string) (*duoapi.StatResult, error) {
	path := fmt.Sprintf("/admin/v1/phones/%s", phoneID)

	_, body, err := c.SignedCall(http.MethodDelete, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &duoapi.StatResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Token methods

// GetTokensTypeAndSerial sets the optional type and serial parameters for a GetTokens request.
func GetTokensTypeAndSerial(typ, serial string) func(*url.Values) {
	return func(opts *url.Values) {
		opts.Set("type", typ)
		opts.Set("serial", serial)
	}
}

// GetTokensResult models responses containing a list of tokens.
type GetTokensResult struct {
	duoapi.StatResult
	ListResult
	Response []Token
}

func (result *GetTokensResult) getResponse() interface{} {
	return result.Response
}

func (result *GetTokensResult) appendResponse(tokens interface{}) {
	asserted_tokens := tokens.([]Token)
	result.Response = append(result.Response, asserted_tokens...)
}

// GetTokens calls GET /admin/v1/tokens
// See https://duo.com/docs/adminapi#retrieve-hardware-tokens
func (c *Client) GetTokens(options ...func(*url.Values)) (*GetTokensResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	cb := func(params url.Values) (responsePage, error) {
		return c.retrieveTokens(params)
	}
	response, err := c.retrieveItems(params, cb)
	if err != nil {
		return nil, err
	}

	return response.(*GetTokensResult), nil
}

func (c *Client) retrieveTokens(params url.Values) (*GetTokensResult, error) {
	_, body, err := c.SignedCall(http.MethodGet, "/admin/v1/tokens", params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetTokensResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetTokenResult models responses containing a single token.
type GetTokenResult struct {
	duoapi.StatResult
	Response Token
}

// GetToken calls GET /admin/v1/tokens/:token_id
// See https://duo.com/docs/adminapi#retrieve-hardware-tokens
func (c *Client) GetToken(tokenID string) (*GetTokenResult, error) {
	path := fmt.Sprintf("/admin/v1/tokens/%s", tokenID)

	_, body, err := c.SignedCall(http.MethodGet, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetTokenResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// U2F token methods

// GetU2FTokensResult models responses containing a list of U2F tokens.
type GetU2FTokensResult struct {
	duoapi.StatResult
	ListResult
	Response []U2FToken
}

func (result *GetU2FTokensResult) getResponse() interface{} {
	return result.Response
}

func (result *GetU2FTokensResult) appendResponse(tokens interface{}) {
	asserted_tokens := tokens.([]U2FToken)
	result.Response = append(result.Response, asserted_tokens...)
}

// GetU2FTokens calls GET /admin/v1/u2ftokens
// See https://duo.com/docs/adminapi#retrieve-u2f-tokens
func (c *Client) GetU2FTokens(options ...func(*url.Values)) (*GetU2FTokensResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	cb := func(params url.Values) (responsePage, error) {
		return c.retrieveU2FTokens(params)
	}
	response, err := c.retrieveItems(params, cb)
	if err != nil {
		return nil, err
	}

	return response.(*GetU2FTokensResult), nil
}

func (c *Client) retrieveU2FTokens(params url.Values) (*GetU2FTokensResult, error) {
	_, body, err := c.SignedCall(http.MethodGet, "/admin/v1/u2ftokens", params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetU2FTokensResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetU2FToken calls GET /admin/v1/u2ftokens/:registration_id
// See https://duo.com/docs/adminapi#retrieve-u2f-token-by-id
func (c *Client) GetU2FToken(registrationID string) (*GetU2FTokensResult, error) {
	path := fmt.Sprintf("/admin/v1/u2ftokens/%s", registrationID)

	_, body, err := c.SignedCall(http.MethodGet, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetU2FTokensResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Integration methods

// GetIntegrationsResult models responses containing a list of integrations.
type GetIntegrationsResult struct {
	duoapi.StatResult
	ListResult
	Response []Integration
}

func (result *GetIntegrationsResult) getResponse() interface{} {
	return result.Response
}

func (result *GetIntegrationsResult) appendResponse(integrations interface{}) {
	asserted_integrations := integrations.([]Integration)
	result.Response = append(result.Response, asserted_integrations...)
}

// GetIntegrations calls GET /admin/v1/integrations
// See https://duo.com/docs/adminapi#retrieve-integrations
func (c *Client) GetIntegrations(options ...func(*url.Values)) (*GetIntegrationsResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	cb := func(params url.Values) (responsePage, error) {
		return c.retrieveIntegrations(params)
	}
	response, err := c.retrieveItems(params, cb)
	if err != nil {
		return nil, err
	}

	return response.(*GetIntegrationsResult), nil
}

func (c *Client) retrieveIntegrations(params url.Values) (*GetIntegrationsResult, error) {
	_, body, err := c.SignedCall(http.MethodGet, "/admin/v1/integrations", params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetIntegrationsResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetIntegrationResult models responses containing a single integration.
type GetIntegrationResult struct {
	duoapi.StatResult
	Response Integration
}

// GetIntegration calls GET /admin/v1/integrations/:integration_key
// See https://duo.com/docs/adminapi#retrieve-integration-by-integration-key
func (c *Client) GetIntegration(integrationKey string) (*GetIntegrationResult, error) {
	path := fmt.Sprintf("/admin/v1/integrations/%s", integrationKey)

	_, body, err := c.SignedCall(http.MethodGet, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetIntegrationResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Administrator methods

// GetAdministratorsResult models responses containing a list of administrators.
type GetAdministratorsResult struct {
	duoapi.StatResult
	ListResult
	Response []Administrator
}

func (result *GetAdministratorsResult) getResponse() interface{} {
	return result.Response
}

func (result *GetAdministratorsResult) appendResponse(administrators interface{}) {
	asserted_administrators := administrators.([]Administrator)
	result.Response = append(result.Response, asserted_administrators...)
}

// GetAdministrators calls GET /admin/v1/administrators
// See https://duo.com/docs/adminapi#retrieve-administrators
func (c *Client) GetAdministrators(options ...func(*url.Values)) (*GetAdministratorsResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	cb := func(params url.Values) (responsePage, error) {
		return c.retrieveAdministrators(params)
	}
	response, err := c.retrieveItems(params, cb)
	if err != nil {
		return nil, err
	}

	return response.(*GetAdministratorsResult), nil
}

func (c *Client) retrieveAdministrators(params url.Values) (*GetAdministratorsResult, error) {
	_, body, err := c.SignedCall(http.MethodGet, "/admin/v1/admins", params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetAdministratorsResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetAdministratorResult models responses containing a single administrator.
type GetAdministratorResult struct {
	duoapi.StatResult
	Response Administrator
}

// GetAdministrator calls GET /admin/v1/administrators/:admin_id
// See https://duo.com/docs/adminapi#retrieve-administrator-by-id
func (c *Client) GetAdministrator(administratorID string) (*GetAdministratorResult, error) {
	path := fmt.Sprintf("/admin/v1/admins/%s", administratorID)

	_, body, err := c.SignedCall(http.MethodGet, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetAdministratorResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// AdministrativeUnit methods

// GetAdministrativeUnitsResult models responses containing a list of administrativeUnits.
type GetAdministrativeUnitsResult struct {
	duoapi.StatResult
	ListResult
	Response []AdministrativeUnit
}

func (result *GetAdministrativeUnitsResult) getResponse() interface{} {
	return result.Response
}

func (result *GetAdministrativeUnitsResult) appendResponse(administrativeUnits interface{}) {
	asserted_administrative_units := administrativeUnits.([]AdministrativeUnit)
	result.Response = append(result.Response, asserted_administrative_units...)
}

// GetAdministrativeUnits calls GET /admin/v1/administrative_units
// See https://duo.com/docs/adminapi#retrieve-administrative-units
func (c *Client) GetAdministrativeUnits(options ...func(*url.Values)) (*GetAdministrativeUnitsResult, error) {
	params := url.Values{}
	for _, o := range options {
		o(&params)
	}

	cb := func(params url.Values) (responsePage, error) {
		return c.retrieveAdministrativeUnits(params)
	}
	response, err := c.retrieveItems(params, cb)
	if err != nil {
		return nil, err
	}

	return response.(*GetAdministrativeUnitsResult), nil
}

func (c *Client) retrieveAdministrativeUnits(params url.Values) (*GetAdministrativeUnitsResult, error) {
	_, body, err := c.SignedCall(http.MethodGet, "/admin/v1/administrative_units", params, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetAdministrativeUnitsResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetAdministrativeUnitResult models responses containing a single administrative unit.
type GetAdministrativeUnitResult struct {
	duoapi.StatResult
	Response AdministrativeUnit
}

// GetAdministrativeUnits calls GET /admin/v1/administrative_units/[admin_unit_id]
// See https://duo.com/docs/adminapi#retrieve-administrative-unit-details
func (c *Client) GetAdministrativeUnit(administrativeUnitID string) (*GetAdministrativeUnitResult, error) {
	path := fmt.Sprintf("/admin/v1/administrative_units/%s", administrativeUnitID)

	_, body, err := c.SignedCall(http.MethodGet, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetAdministrativeUnitResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Account Info methods

type AccountInfo struct {
	AdminCount                int `json:"admin_count"`
	IntegrationCount          int `json:"integration_count"`
	TelephonyCreditsRemaining int `json:"telephony_credits_remaining"`
	UserCount                 int `json:"user_count"`
	UserPendingDeletionCount  int `json:"user_pending_deletion_count"`
}

// GetAccountInfoSummaryResult model responses for the account info summary.
type GetAccountInfoSummaryResult struct {
	duoapi.StatResult
	Response AccountInfo
}

// GetAccountInfoSummary calls GET /admin/v1/info/summary
// See https://duo.com/docs/adminapi#retrieve-summary
func (c *Client) GetAccountInfoSummary() (*GetAccountInfoSummaryResult, error) {
	path := fmt.Sprintf("/admin/v1/info/summary")

	_, body, err := c.SignedCall(http.MethodGet, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetAccountInfoSummaryResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Account Settings methods

type AccountSettings struct {
	CallerID                   string  `json:"caller_id"`
	FraudEmail                 string  `json:"fraud_email"`
	FraudEmailEnabled          bool    `json:"fraud_email_enabled"`
	HelpdeskBypass             string  `json:"helpdesk_bypass"`
	HelpdeskBypassExpiration   int     `json:"helpdesk_bypass_expiration"`
	HelpdeskCanSendEnrollEmail bool    `json:"helpdesk_can_send_enroll_email"`
	HelpdeskMessage            string  `json:"helpdesk_message"`
	InactiveUserExpiration     int     `json:"inactive_user_expiration"`
	KeypressConfirm            string  `json:"keypress_confirm"`
	KeypressFraud              string  `json:"keypress_fraud"`
	Language                   string  `json:"language"`
	LockoutExpireDuration      int     `json:"lockout_expire_duration"`
	LockoutThreshold           int     `json:"lockout_threshold"`
	MinimumPasswordLength      int     `json:"minimum_password_length"`
	Name                       string  `json:"name"`
	PasswordRequiresLowerAlpha bool    `json:"password_requires_lower_alpha"`
	PasswordRequiresNumeric    bool    `json:"password_requires_numeric"`
	PasswordRequiresSpecial    bool    `json:"password_requires_special"`
	PasswordRequiresUpperAlpha bool    `json:"password_requires_upper_alpha"`
	SmsBatch                   int     `json:"sms_batch"`
	SmsExpiration              int     `json:"sms_expiration"`
	SmsMessage                 string  `json:"sms_message"`
	SmsRefresh                 int     `json:"sms_refresh"`
	TelephonyWarningMin        int     `json:"telephony_warning_min"`
	Timezone                   string  `json:"timezone"`
	UserTelephonyCostMax       float64 `json:"user_telephony_cost_max"`
}

// GetAccountInfoSummaryResult model responses for the account info summary.
type GetAccountSettingsResult struct {
	duoapi.StatResult
	Response AccountSettings
}

// GetAccountInfoSummary calls GET /admin/v1/info/summary
// See https://duo.com/docs/adminapi#retrieve-summary
func (c *Client) GetAccountSettings() (*GetAccountSettingsResult, error) {
	path := fmt.Sprintf("/admin/v1/settings")

	_, body, err := c.SignedCall(http.MethodGet, path, nil, duoapi.UseTimeout)
	if err != nil {
		return nil, err
	}

	result := &GetAccountSettingsResult{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
