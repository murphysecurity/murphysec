package gradle

// language=Groovy
var initScriptContent = `
import groovy.json.JsonOutput


def mergedDepsConfExecuted = false
allprojects { everyProj ->
    task getDepsJson {
        def onlyConf = project.hasProperty('configuration') ? configuration : null

        def subDepsToDict
        subDepsToDict = { deps, deepLevel ->
            def res = [:]
            deepLevel = deepLevel + 1
            if (deepLevel < 4) {
                deps.each{ d ->
                    def row = ['name': "$d.moduleGroup:$d.moduleName", 'version': d.moduleVersion]
                    def subDeps = subDepsToDict(d.children, deepLevel)
					row['dependencies'] = subDeps
                    res[row['name']] = row
                }
            }
            return res

        }

        def depsToDict
        depsToDict = { deps ->  
            def res = [:]
            deps.each { d ->
                def row = ['name': "$d.moduleGroup:$d.moduleName", 'version': d.moduleVersion]
                def deepLevel = 0
                def subDeps = subDepsToDict(d.children,deepLevel)
				row['dependencies'] = subDeps
                res[row['name']] = row
            }
            return res
        }

        doLast { task ->
            def projectsDict = [:]
            def result = ['defaultProject': task.project.name, 'projects': projectsDict]
            if (!mergedDepsConfExecuted) {
                allprojects.each { proj ->
                    def projectConf = null
                    if (proj.configurations.size() > 0) {
                        if (onlyConf != null) {
                            projectConf = proj.configurations.getByName(onlyConf)
                        } else if (proj.configurations.findAll({ it.name == 'mergedDepsConf'}).size() == 0) {
                            projectConf = proj.configurations.create('mergedDepsConf')
                            proj.configurations
                                .findAll({ it.name != 'mergedDepsConf' && (onlyConf == null || it.name == onlyConf) })
                                .each { projectConf.extendsFrom(it) }
                        }
                    }
                    if (projectConf != null) {
                        projectsDict[proj.name] = [
                            'targetFile': findProject(proj.path).buildFile.toString(),
                            'depDict': depsToDict(projectConf.resolvedConfiguration.firstLevelModuleDependencies)
                        ]
                    } else {
                        projectsDict[proj.name] = [
                            'targetFile': findProject(proj.path).buildFile.toString()
                        ]
                    }
                }
                println("GetDepsJson:" + JsonOutput.toJson(result))
                mergedDepsConfExecuted = true
            }
        }
    }
}

`
