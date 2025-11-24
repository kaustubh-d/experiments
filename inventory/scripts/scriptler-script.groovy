import org.yaml.snakeyaml.Yaml

try {
    // ---------- SAFE REST CALL ----------
    def url = new URL("https://raw.githubusercontent.com/kaustubh-d/experiments/refs/heads/app-inventory/inventory/enabled-app-list.yaml")
    //def conn = url.openConnection()
    //conn.setRequestMethod("GET")
    //conn.connect()

    //def loggedInUser = "12345"

    // Read Yaml response
    //def yamlText = conn.inputStream.text
  	def yamlText = url.getText([connectTimeout: 5000, readTimeout: 10000])

    def list = []

    try {
        // Load the entire YAML file into a Map
        // Structure: Map<String, Object> where Object is another Map
        Map<String, Object> allAppsData = new Yaml().load(yamlText);

        // Iterate through each entry in the top-level map (e.g., appname-1, appname-2)
        for (Map.Entry<String, Object> entry : allAppsData.entrySet()) {
            String appName = entry.getKey();
            // The value for each entry is a nested map (the owner data)
            Map<String, Object> ownerData = (LinkedHashMap<String, Object>) entry.getValue();

            // Access the "owner" key within the nested map
            // The value associated with "owner" is a List of Strings/Longs
            List<String> owners = (List<String>) ownerData.get("owner"); // Use String or Long depending on ID format

            // System.out.println("Application: " + appName);
            // System.out.println("Owners: " + owners);
            
            // You can also iterate through the owners list individually
            for (String ownerId : owners) {
                // System.out.println("  Owner ID: " + ownerId);
                if (ownerId == loggedInUser) {
                    list.add(appName)
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