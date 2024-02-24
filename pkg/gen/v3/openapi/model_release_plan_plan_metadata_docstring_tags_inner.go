/*
 * Puppet Forge v3 API
 *
 * ## Introduction The Puppet Forge API (hereafter referred to as the Forge API) provides quick access to all the data on the Puppet Forge via a RESTful interface. Using the Forge API, you can write scripts and tools that interact with the Puppet Forge website.  The Forge API's current version is `v3`. It is considered regression-stable, meaning that the returned data is guaranteed to include the fields described in the schemas on this page; however, additional data might be added in the future and clients must ignore any properties they do not recognize.  ## OpenAPI Specification The Puppet Forge v3 API is described by an [OpenAPI 3.0](https://github.com/OAI/OpenAPI-Specification/blob/main/versions/3.0.0.md) formatted specification file. The most up-to-date version of this specification file can be accessed at [https://forgeapi.puppet.com/v3/openapi.json](/v3/openapi.json).  ## Features * The API is accessed over HTTPS via either the `forgeapi.puppet.com` (IPv4 or IPv6). All data is returned in JSON   format. * Blank fields are included as `null`. * Nested resources may use an abbreviated representation. A link to the full representation for the   resource is always included. * All timestamps in JSON responses are returned in ISO 8601 format: `YYYY-MM-DD HH:MM:SS ±HHMM`. * The HTTP response headers include caching hints for conditional requests.  ## Concepts and Terminology * **Module**: Modules are self-contained bundles of code and data with a specific directory structure. Modules are identified by a combination of the author's username and the module's name, separated by a hyphen. For example: `puppetlabs-apache` * **Release**: A single, specific version of a module is called a Release. Releases are identified by a combination of the module identifier (see above) and the Release version, separated by a hyphen. For example: `puppetlabs-apache-4.0.0`  ## Errors The Puppet Forge API follows [RFC 2616](https://tools.ietf.org/html/rfc2616) and [RFC 6585](https://tools.ietf.org/html/rfc6585).  Error responses are served with a `4xx` or `5xx` status code, and are sent as a JSON document with a content type of `application/json`. The error document contains the following top-level keys and values:    * `message`: a string value that summarizes the problem   * `errors`: a list (array) of strings containing additional details describing the underlying cause(s) of the     failure  An example error response is shown below:  ```json {   \"message\": \"400 Bad Request\",   \"errors\": [     \"Cannot parse request body as JSON\"   ] } ```  ## User-Agent Required All API requests must include a valid `User-Agent` header. Requests with no `User-Agent` header will be rejected. The `User-Agent` header helps identify your application or library, so we can communicate with you if necessary. If your use of the API is informal or personal, we recommend using your username as the value for the `User-Agent` header.  User-Agent headers are a list of one or more product descriptions, generally taking this form:  ``` <name-without-spaces>/<version> (comments) ```  For example, the following are all useful User-Agent values:  ``` MyApplication/0.0.0 Her/0.6.8 Faraday/0.8.8 Ruby/1.9.3-p194 (i386-linux) My-Library-Name/1.2.4 myusername ```  ## Hostname Configuration Most tools that interact with the Forge API allow specification of the hostname to use. You can configure a few common tools to use a specified hostname as follows:  For **Puppet Enterprise** users, in [r10k](https://puppet.com/docs/pe/latest/r10k_customize_config.html#r10k_configuring_forge_settings) or [Code Manager](https://puppet.com/docs/pe/latest/code_mgr_customizing.html#config_forge_settings), specify `forge_settings` in Hiera: ``` pe_r10k::forge_settings:   baseurl: 'https://forgeapi.puppet.com' ``` or ``` puppet_enterprise::master::code_manager::forge_settings:   baseurl: 'https://forgeapi.puppet.com' ``` <br />  If you are an **open source Puppet** user using r10k, you'll need to [edit your r10k.yaml directly](https://github.com/puppetlabs/r10k/blob/main/doc/dynamic-environments/configuration.mkd#forge): ``` forge:   baseurl: 'https://forgeapi.puppet.com' ``` or set the appropriate class param for the [open source r10k module](https://forge.puppet.com/puppet/r10k#forge_settings): ``` $forge_settings = {   'baseurl' => 'https://forgeapi.puppet.com', } ``` <br />  In [**Bolt**](https://puppet.com/docs/bolt/latest/bolt_installing_modules.html#install-forge-modules-from-an-alternate-forge), set a `baseurl` for the Forge in `bolt-project.yaml`: ``` module-install:   forge:     baseurl: https://forgeapi.puppet.com ``` <br />  Using `puppet config`: ``` $ puppet config set module_repository https://forgeapi.puppet.com ``` 
 *
 * API version: 29
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package gorge




// ReleasePlanPlanMetadataDocstringTagsInner - Generated docstring for a specific YARD tag (like @param, for example)
type ReleasePlanPlanMetadataDocstringTagsInner struct {

	// The name specified by this YARD tag
	Name string `json:"name,omitempty"`

	// The description text specified by this YARD tag
	Text string `json:"text,omitempty"`

	// Valid types specified by this YARD tag
	Types []string `json:"types,omitempty"`

	// The name of the YARD tag itself
	TagName string `json:"tag_name,omitempty"`
}

// AssertReleasePlanPlanMetadataDocstringTagsInnerRequired checks if the required fields are not zero-ed
func AssertReleasePlanPlanMetadataDocstringTagsInnerRequired(obj ReleasePlanPlanMetadataDocstringTagsInner) error {
	return nil
}

// AssertReleasePlanPlanMetadataDocstringTagsInnerConstraints checks if the values respects the defined constraints
func AssertReleasePlanPlanMetadataDocstringTagsInnerConstraints(obj ReleasePlanPlanMetadataDocstringTagsInner) error {
	return nil
}
