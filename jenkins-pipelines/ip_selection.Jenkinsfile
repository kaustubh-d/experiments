@Library('custom-helpers') _

pipeline {
    agent any

    parameters {
        hidden(defaultValue: 'host1:10.0.0.1,host2:10.0.0.2,host3:10.0.0.3,host4:10.0.0.4', name: 'prod_ip_list')

        booleanParam(description: 'Select/Unselect all', name: 'select_all_ips_flag', defaultValue: false)

        // Dynamic choice rendered in Jenkins UI using Groovy script
        reactiveChoice(
            name: 'selected_hostips',
            choiceType: 'PT_CHECKBOX',
            filterable: true,
            referencedParameters: 'prod_ip_list,select_all_ips_flag', defaultValue: false,
            // ... other settings ...
            script: [
                $class: 'GroovyScript',
                script: [
                    script: '''
                        def csvStrToList(str) {
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
                        def allIps = csvStrToList(prod_ip_list)

                        def selectedIpList = []
                        def suffix = ""
                        if (select_all_ips_flag) {
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
                def selectedIpList = []
                def selectedHostIpList = stringHelpers.csvStrToList(params.selected_hostips)
                selectedHostIpList.each { ip ->
                    def extractedIp = stringHelpers.ipFromHostIp(ip)
                    selectedIpList.add(extractedIp)
                }
                echo "Selected IPs: ${selectedIpList}"
                echo "select_all_ips_flag: ${params.select_all_ips_flag}"
                echo "selected_hostips: ${params.selected_hostips}"
            }
          }
        }
    }
}