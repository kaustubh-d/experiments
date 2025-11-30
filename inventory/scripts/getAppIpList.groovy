import groovy.json.JsonSlurper

// Usage: getAppIpList("appName", "envName") -> List of IP addresses for the specified app and environment

def getAppIpList(String appName, envName, loggedInUser) {
	try {
    // ---------- SAFE REST CALL ----------
    url_path = "http://localhost:9080/inventory/apps/${appName}/${envName}"

    def url = new URL(url_path)

    def jsonText = url.getText([connectTimeout: 5000, readTimeout: 10000])

    // Convert Json to groovy object
    def jsonSlurper = new groovy.json.JsonSlurper()
    def object = jsonSlurper.parseText(jsonText)

    def list = []

    try {
        // Assume response is a single app object at the root
        if (object instanceof Map) {
          // Get access list and AppEnvData
          def access = object.get("access")
          def appEnvData = object.get("AppEnvData") ?: object.get("appEnvData") ?: object.get("AppEnv")
          if (access) {
            def accessList = (access instanceof Collection) ? access.collect { it?.toString() } : [access.toString()]
            if (accessList.contains(loggedInUser?.toString())) {
              def hosts = appEnvData?.get("hosts")
              if (hosts instanceof Collection) {
                hosts.each { host ->
                  if (host instanceof Map && host.ip_address) {
                    list.add(host.ip_address.toString())
                  }
                }
              }
            }
          }
        }
    } catch (Exception e) {
        return ["Error parsing YAML: ${e.message}"]
    }

    return list
  }
  catch (Exception e) {
      // Fail-safe fallback to avoid UI error
      return ["Error fetching data: ${e.message}"]
  }
  return ["Error fetching data"]
}

return this