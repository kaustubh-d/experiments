pipeline {
    agent any

    parameters {
        hidden(defaultValue: 'host1:10.0.0.1,host2:10.0.0.2,host3:10.0.0.3,host4:10.0.0.4', name: 'prod_ip_list')

        booleanParam(description: 'Select/Unselect all', name: 'select_all_ips')

        // Dynamic choice rendered in Jenkins UI using Groovy script
        reactiveChoice(
            name: 'selected_ips',
            choiceType: 'PT_CHECKBOX',
            filterable: true,
            referencedParameters: 'prod_ip_list,select_all_ips',
            // ... other settings ...
            script: [
                $class: 'GroovyScript',
                script: [
                    script: '''
                        def strToList(str) {
                            if (!str) {
                                return []
                            }
                            return str.split(',').collect { it.trim() }
                        }
                        def assertNotEmpty(tag, str) {
                            if (!str || str.trim().length() == 0) {
                                throw new Exception("${tag} is empty")
                            }
                        }
                        assertNotEmpty("prod_ip_list", prod_ip_list)
                        def allIps = strToList(prod_ip_list)

                        def selectedIpList = []
                        def suffix = ""
                        if (select_all_ips) {
                            suffix = ":selected"
                        }
                        allIps.each { ip ->
                            selectedIpList.add(ip + suffix)
                        }
                        return selectedIpList
                    '''.stripIndent(),
                    sandbox: true
                ]
            ]
        )
    }

    stages {
        stage('Show Selected') {
            steps {
              script {
                def ipFromHostIp = { hostIp ->
                    if (!hostIp) {
                        return null
                    }
                    def parts = hostIp.split(':')
                    if (parts.length == 2) {
                        return parts[1].trim()
                    }
                    return null
                }
                def selectedIpList = []
                def selectedHostIpList = (params.selected_ips ?: '').split(',').collect { it.trim() }.findAll { it }
                selectedHostIpList.each { ip ->
                    def extractedIp = ipFromHostIp(ip)
                    selectedIpList.add(extractedIp)
                }
                echo "selectedIpList: ${selectedIpList.findAll{ it }.join(', ')}"
                echo "select_all_ips: ${params.select_all_ips}"
                echo "selected_ips: ${params.selected_ips}"
            }
          }
        }
    }
}