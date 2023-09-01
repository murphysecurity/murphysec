package maven

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vifraa/gopom"
	"testing"
)

func TestDependencySet(t *testing.T) {
	// language=json
	var ref = `
       [
         {
           "GroupID": "org.projectlombok",
           "ArtifactID": "lombok",
           "version": "1.18.8",
           "Type": "",
           "Classifier": "",
           "Scope": "provided",
           "SystemPath": "",
           "Exclusions": null,
           "Optional": ""
         },
         {
           "GroupID": "com.monitorjbl",
           "ArtifactID": "xlsx-streamer",
           "version": "2.2.0",
           "Type": "",
           "Classifier": "",
           "Scope": "",
           "SystemPath": "",
           "Exclusions": [
             {
               "ArtifactID": "poi",
               "GroupID": "org.apache.poi"
             },
             {
               "ArtifactID": "poi-ooxml",
               "GroupID": "org.apache.poi"
             },
             {
               "ArtifactID": "poi-ooxml-schemas",
               "GroupID": "org.apache.poi"
             }
           ],
           "Optional": ""
         },
         {
           "GroupID": "org.apache.poi",
           "ArtifactID": "poi-ooxml",
           "version": "4.1.2",
           "Type": "",
           "Classifier": "",
           "Scope": "",
           "SystemPath": "",
           "Exclusions": [
             {
               "ArtifactID": "commons-compress",
               "GroupID": "org.apache.commons"
             }
           ],
           "Optional": ""
         },
         {
           "GroupID": "org.apache.poi",
           "ArtifactID": "poi-ooxml-schemas",
           "version": "4.1.2",
           "Type": "",
           "Classifier": "",
           "Scope": "",
           "SystemPath": "",
           "Exclusions": null,
           "Optional": ""
         }
       ]`
	// language=xml
	var p1 = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
   <properties>
       <poi.version>4.1.2</poi.version>
       <jackson.version>2.13.0</jackson.version>
   </properties>
   <dependencies>
       <dependency>
           <groupId>org.projectlombok</groupId>
           <artifactId>lombok</artifactId>
       </dependency>
       <dependency>
           <groupId>com.monitorjbl</groupId>
           <artifactId>xlsx-streamer</artifactId>
           <version>2.2.0</version>
           <exclusions>
               <exclusion>
                   <groupId>org.apache.poi</groupId>
                   <artifactId>poi-ooxml-schemas</artifactId>
               </exclusion>
               <exclusion>
                   <groupId>org.apache.poi</groupId>
                   <artifactId>poi</artifactId>
               </exclusion>
           </exclusions>
       </dependency>
   </dependencies>
</project>`
	// language=xml
	var p2 = `<?xml version="1.0" encoding="UTF-8"?>
<project xmlns="http://maven.apache.org/POM/4.0.0" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:schemaLocation="http://maven.apache.org/POM/4.0.0 http://maven.apache.org/xsd/maven-4.0.0.xsd">
   <properties>
       <poi.version>4.1.2</poi.version>
       <jackson.version>2.13.0</jackson.version>
   </properties>
   <dependencyManagement>
       <dependencies>
           <dependency>
               <groupId>org.projectlombok</groupId>
               <artifactId>lombok</artifactId>
               <version>1.18.8</version>
               <scope>provided</scope>
               <optional>true</optional>
           </dependency>
       </dependencies>
   </dependencyManagement>
   <dependencies>
       <dependency>
           <groupId>org.apache.poi</groupId>
           <artifactId>poi-ooxml</artifactId>
           <version>${poi.version}</version>
           <exclusions>
               <exclusion>
                   <groupId>org.apache.commons</groupId>
                   <artifactId>commons-compress</artifactId>
               </exclusion>
           </exclusions>
       </dependency>
       <dependency>
           <groupId>org.apache.poi</groupId>
           <artifactId>poi-ooxml-schemas</artifactId>
           <version>${poi.version}</version>
       </dependency>
       <dependency>
           <groupId>com.monitorjbl</groupId>
           <artifactId>xlsx-streamer</artifactId>
           <version>2.2.0</version>
           <exclusions>
               <exclusion>
                   <groupId>org.apache.poi</groupId>
                   <artifactId>poi-ooxml</artifactId>
               </exclusion>
           </exclusions>
       </dependency>
   </dependencies>
</project>`

	var x1, x2 gopom.Project
	assert.NoError(t, xml.Unmarshal([]byte(p1), &x1))
	assert.NoError(t, xml.Unmarshal([]byte(p2), &x2))

	var dp = newPomDependencySet()
	dp.mergeAll(x1.Dependencies, true, false)
	dp.mergeAll(x2.Dependencies, true, false)
	var ip = newProperties()
	ip.PutMap(x1.Properties.Entries)
	ip.PutMap(x2.Properties.Entries)
	dp.mergeProperty(ip)
	var dm = newPomDependencySet()
	dm.mergeAll(x1.DependencyManagement.Dependencies, true, false)
	dm.mergeAll(x2.DependencyManagement.Dependencies, true, false)
	dp.mergeAll(dm.listAll(), false, true)

	var expect []gopom.Dependency
	assert.NoError(t, json.Unmarshal([]byte(ref), &expect))
	var expectTexts []string
	for _, it := range expect {
		expectTexts = append(expectTexts, fmt.Sprint(it))
	}
	var actualTexts []string
	for _, it := range dp.listAll() {
		actualTexts = append(actualTexts, fmt.Sprint(it))
	}
	t.Log(actualTexts)
	assert.EqualValues(t, expectTexts, actualTexts)
}
