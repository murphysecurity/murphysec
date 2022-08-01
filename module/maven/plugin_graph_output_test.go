package maven

import (
	"github.com/murphysecurity/murphysec/utils/must"
	"os"
	"testing"
)

func TestDependencyGraph_ReadFromFile(t *testing.T) {
	// language=json
	var a = `{
  "graphName" : "mall",
  "artifacts" : [ {
    "id" : "org.springframework.boot:spring-boot-starter:jar",
    "numericId" : 1,
    "groupId" : "org.springframework.boot",
    "artifactId" : "spring-boot-starter",
    "version" : "2.3.0.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework.boot:spring-boot:jar",
    "numericId" : 2,
    "groupId" : "org.springframework.boot",
    "artifactId" : "spring-boot",
    "version" : "2.3.0.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework.boot:spring-boot-autoconfigure:jar",
    "numericId" : 3,
    "groupId" : "org.springframework.boot",
    "artifactId" : "spring-boot-autoconfigure",
    "version" : "2.3.0.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "ch.qos.logback:logback-classic:jar",
    "numericId" : 4,
    "groupId" : "ch.qos.logback",
    "artifactId" : "logback-classic",
    "version" : "1.2.3",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "ch.qos.logback:logback-core:jar",
    "numericId" : 5,
    "groupId" : "ch.qos.logback",
    "artifactId" : "logback-core",
    "version" : "1.2.3",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework.boot:spring-boot-starter-logging:jar",
    "numericId" : 6,
    "groupId" : "org.springframework.boot",
    "artifactId" : "spring-boot-starter-logging",
    "version" : "2.3.0.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.apache.logging.log4j:log4j-to-slf4j:jar",
    "numericId" : 7,
    "groupId" : "org.apache.logging.log4j",
    "artifactId" : "log4j-to-slf4j",
    "version" : "2.13.2",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.apache.logging.log4j:log4j-api:jar",
    "numericId" : 8,
    "groupId" : "org.apache.logging.log4j",
    "artifactId" : "log4j-api",
    "version" : "2.13.2",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.slf4j:jul-to-slf4j:jar",
    "numericId" : 9,
    "groupId" : "org.slf4j",
    "artifactId" : "jul-to-slf4j",
    "version" : "1.7.30",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "jakarta.annotation:jakarta.annotation-api:jar",
    "numericId" : 10,
    "groupId" : "jakarta.annotation",
    "artifactId" : "jakarta.annotation-api",
    "version" : "1.3.5",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.yaml:snakeyaml:jar",
    "numericId" : 11,
    "groupId" : "org.yaml",
    "artifactId" : "snakeyaml",
    "version" : "1.26",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework.boot:spring-boot-starter-actuator:jar",
    "numericId" : 12,
    "groupId" : "org.springframework.boot",
    "artifactId" : "spring-boot-starter-actuator",
    "version" : "2.3.0.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework.boot:spring-boot-actuator-autoconfigure:jar",
    "numericId" : 13,
    "groupId" : "org.springframework.boot",
    "artifactId" : "spring-boot-actuator-autoconfigure",
    "version" : "2.3.0.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework.boot:spring-boot-actuator:jar",
    "numericId" : 14,
    "groupId" : "org.springframework.boot",
    "artifactId" : "spring-boot-actuator",
    "version" : "2.3.0.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "com.fasterxml.jackson.core:jackson-databind:jar",
    "numericId" : 15,
    "groupId" : "com.fasterxml.jackson.core",
    "artifactId" : "jackson-databind",
    "version" : "2.11.0",
    "optional" : false,
    "scopes" : [ "runtime" ],
    "types" : [ "jar" ]
  }, {
    "id" : "com.fasterxml.jackson.core:jackson-annotations:jar",
    "numericId" : 16,
    "groupId" : "com.fasterxml.jackson.core",
    "artifactId" : "jackson-annotations",
    "version" : "2.11.0",
    "optional" : false,
    "scopes" : [ "runtime" ],
    "types" : [ "jar" ]
  }, {
    "id" : "com.fasterxml.jackson.core:jackson-core:jar",
    "numericId" : 17,
    "groupId" : "com.fasterxml.jackson.core",
    "artifactId" : "jackson-core",
    "version" : "2.11.0",
    "optional" : false,
    "scopes" : [ "runtime" ],
    "types" : [ "jar" ]
  }, {
    "id" : "com.fasterxml.jackson.datatype:jackson-datatype-jsr310:jar",
    "numericId" : 18,
    "groupId" : "com.fasterxml.jackson.datatype",
    "artifactId" : "jackson-datatype-jsr310",
    "version" : "2.11.0",
    "optional" : false,
    "scopes" : [ "runtime" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework:spring-context:jar",
    "numericId" : 19,
    "groupId" : "org.springframework",
    "artifactId" : "spring-context",
    "version" : "5.2.6.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework:spring-expression:jar",
    "numericId" : 20,
    "groupId" : "org.springframework",
    "artifactId" : "spring-expression",
    "version" : "5.2.6.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "io.micrometer:micrometer-core:jar",
    "numericId" : 21,
    "groupId" : "io.micrometer",
    "artifactId" : "micrometer-core",
    "version" : "1.5.1",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.hdrhistogram:HdrHistogram:jar",
    "numericId" : 22,
    "groupId" : "org.hdrhistogram",
    "artifactId" : "HdrHistogram",
    "version" : "2.1.12",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.latencyutils:LatencyUtils:jar",
    "numericId" : 23,
    "groupId" : "org.latencyutils",
    "artifactId" : "LatencyUtils",
    "version" : "2.0.3",
    "optional" : false,
    "scopes" : [ "runtime" ],
    "types" : [ "jar" ]
  }, {
    "id" : "com.macro.mall:mall:pom",
    "numericId" : 24,
    "groupId" : "com.macro.mall",
    "artifactId" : "mall",
    "version" : "1.0-SNAPSHOT",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "pom" ]
  }, {
    "id" : "org.springframework:spring-aop:jar",
    "numericId" : 25,
    "groupId" : "org.springframework",
    "artifactId" : "spring-aop",
    "version" : "5.2.6.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework:spring-beans:jar",
    "numericId" : 26,
    "groupId" : "org.springframework",
    "artifactId" : "spring-beans",
    "version" : "5.2.6.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework.boot:spring-boot-starter-aop:jar",
    "numericId" : 27,
    "groupId" : "org.springframework.boot",
    "artifactId" : "spring-boot-starter-aop",
    "version" : "2.3.0.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.aspectj:aspectjweaver:jar",
    "numericId" : 28,
    "groupId" : "org.aspectj",
    "artifactId" : "aspectjweaver",
    "version" : "1.9.5",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework.boot:spring-boot-starter-test:jar",
    "numericId" : 29,
    "groupId" : "org.springframework.boot",
    "artifactId" : "spring-boot-starter-test",
    "version" : "2.3.0.RELEASE",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework.boot:spring-boot-test:jar",
    "numericId" : 30,
    "groupId" : "org.springframework.boot",
    "artifactId" : "spring-boot-test",
    "version" : "2.3.0.RELEASE",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework.boot:spring-boot-test-autoconfigure:jar",
    "numericId" : 31,
    "groupId" : "org.springframework.boot",
    "artifactId" : "spring-boot-test-autoconfigure",
    "version" : "2.3.0.RELEASE",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "net.minidev:accessors-smart:jar",
    "numericId" : 32,
    "groupId" : "net.minidev",
    "artifactId" : "accessors-smart",
    "version" : "1.2",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.ow2.asm:asm:jar",
    "numericId" : 33,
    "groupId" : "org.ow2.asm",
    "artifactId" : "asm",
    "version" : "5.0.4",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "net.minidev:json-smart:jar",
    "numericId" : 34,
    "groupId" : "net.minidev",
    "artifactId" : "json-smart",
    "version" : "2.3",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "com.jayway.jsonpath:json-path:jar",
    "numericId" : 35,
    "groupId" : "com.jayway.jsonpath",
    "artifactId" : "json-path",
    "version" : "2.4.0",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.slf4j:slf4j-api:jar",
    "numericId" : 36,
    "groupId" : "org.slf4j",
    "artifactId" : "slf4j-api",
    "version" : "1.7.30",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "jakarta.xml.bind:jakarta.xml.bind-api:jar",
    "numericId" : 37,
    "groupId" : "jakarta.xml.bind",
    "artifactId" : "jakarta.xml.bind-api",
    "version" : "2.3.3",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "jakarta.activation:jakarta.activation-api:jar",
    "numericId" : 38,
    "groupId" : "jakarta.activation",
    "artifactId" : "jakarta.activation-api",
    "version" : "1.2.2",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.assertj:assertj-core:jar",
    "numericId" : 39,
    "groupId" : "org.assertj",
    "artifactId" : "assertj-core",
    "version" : "3.16.1",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.hamcrest:hamcrest:jar",
    "numericId" : 40,
    "groupId" : "org.hamcrest",
    "artifactId" : "hamcrest",
    "version" : "2.2",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.junit.jupiter:junit-jupiter-api:jar",
    "numericId" : 41,
    "groupId" : "org.junit.jupiter",
    "artifactId" : "junit-jupiter-api",
    "version" : "5.6.2",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.opentest4j:opentest4j:jar",
    "numericId" : 42,
    "groupId" : "org.opentest4j",
    "artifactId" : "opentest4j",
    "version" : "1.2.0",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.junit.platform:junit-platform-commons:jar",
    "numericId" : 43,
    "groupId" : "org.junit.platform",
    "artifactId" : "junit-platform-commons",
    "version" : "1.6.2",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.junit.jupiter:junit-jupiter:jar",
    "numericId" : 44,
    "groupId" : "org.junit.jupiter",
    "artifactId" : "junit-jupiter",
    "version" : "5.6.2",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.junit.jupiter:junit-jupiter-params:jar",
    "numericId" : 45,
    "groupId" : "org.junit.jupiter",
    "artifactId" : "junit-jupiter-params",
    "version" : "5.6.2",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.junit.jupiter:junit-jupiter-engine:jar",
    "numericId" : 46,
    "groupId" : "org.junit.jupiter",
    "artifactId" : "junit-jupiter-engine",
    "version" : "5.6.2",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.junit.vintage:junit-vintage-engine:jar",
    "numericId" : 47,
    "groupId" : "org.junit.vintage",
    "artifactId" : "junit-vintage-engine",
    "version" : "5.6.2",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.apiguardian:apiguardian-api:jar",
    "numericId" : 48,
    "groupId" : "org.apiguardian",
    "artifactId" : "apiguardian-api",
    "version" : "1.1.0",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.junit.platform:junit-platform-engine:jar",
    "numericId" : 49,
    "groupId" : "org.junit.platform",
    "artifactId" : "junit-platform-engine",
    "version" : "1.6.2",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "junit:junit:jar",
    "numericId" : 50,
    "groupId" : "junit",
    "artifactId" : "junit",
    "version" : "4.13",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.mockito:mockito-core:jar",
    "numericId" : 51,
    "groupId" : "org.mockito",
    "artifactId" : "mockito-core",
    "version" : "3.3.3",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "net.bytebuddy:byte-buddy:jar",
    "numericId" : 52,
    "groupId" : "net.bytebuddy",
    "artifactId" : "byte-buddy",
    "version" : "1.10.10",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "net.bytebuddy:byte-buddy-agent:jar",
    "numericId" : 53,
    "groupId" : "net.bytebuddy",
    "artifactId" : "byte-buddy-agent",
    "version" : "1.10.10",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.objenesis:objenesis:jar",
    "numericId" : 54,
    "groupId" : "org.objenesis",
    "artifactId" : "objenesis",
    "version" : "2.6",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.mockito:mockito-junit-jupiter:jar",
    "numericId" : 55,
    "groupId" : "org.mockito",
    "artifactId" : "mockito-junit-jupiter",
    "version" : "3.3.3",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.skyscreamer:jsonassert:jar",
    "numericId" : 56,
    "groupId" : "org.skyscreamer",
    "artifactId" : "jsonassert",
    "version" : "1.5.0",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "com.vaadin.external.google:android-json:jar",
    "numericId" : 57,
    "groupId" : "com.vaadin.external.google",
    "artifactId" : "android-json",
    "version" : "0.0.20131108.vaadin1",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework:spring-core:jar",
    "numericId" : 58,
    "groupId" : "org.springframework",
    "artifactId" : "spring-core",
    "version" : "5.2.6.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework:spring-jcl:jar",
    "numericId" : 59,
    "groupId" : "org.springframework",
    "artifactId" : "spring-jcl",
    "version" : "5.2.6.RELEASE",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework:spring-test:jar",
    "numericId" : 60,
    "groupId" : "org.springframework",
    "artifactId" : "spring-test",
    "version" : "5.2.6.RELEASE",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.xmlunit:xmlunit-core:jar",
    "numericId" : 61,
    "groupId" : "org.xmlunit",
    "artifactId" : "xmlunit-core",
    "version" : "2.7.0",
    "optional" : false,
    "scopes" : [ "test" ],
    "types" : [ "jar" ]
  }, {
    "id" : "cn.hutool:hutool-all:jar",
    "numericId" : 62,
    "groupId" : "cn.hutool",
    "artifactId" : "hutool-all",
    "version" : "5.4.0",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.projectlombok:lombok:jar",
    "numericId" : 63,
    "groupId" : "org.projectlombok",
    "artifactId" : "lombok",
    "version" : "1.18.12",
    "optional" : false,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  }, {
    "id" : "org.springframework.boot:spring-boot-configuration-processor:jar",
    "numericId" : 64,
    "groupId" : "org.springframework.boot",
    "artifactId" : "spring-boot-configuration-processor",
    "version" : "2.3.0.RELEASE",
    "optional" : true,
    "scopes" : [ "compile" ],
    "types" : [ "jar" ]
  } ],
  "dependencies" : [ {
    "from" : "org.springframework.boot:spring-boot-starter:jar",
    "to" : "org.springframework.boot:spring-boot:jar",
    "numericFrom" : 0,
    "numericTo" : 1,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter:jar",
    "to" : "org.springframework.boot:spring-boot-autoconfigure:jar",
    "numericFrom" : 0,
    "numericTo" : 2,
    "resolution" : "INCLUDED"
  }, {
    "from" : "ch.qos.logback:logback-classic:jar",
    "to" : "ch.qos.logback:logback-core:jar",
    "numericFrom" : 3,
    "numericTo" : 4,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-logging:jar",
    "to" : "ch.qos.logback:logback-classic:jar",
    "numericFrom" : 5,
    "numericTo" : 3,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.apache.logging.log4j:log4j-to-slf4j:jar",
    "to" : "org.apache.logging.log4j:log4j-api:jar",
    "numericFrom" : 6,
    "numericTo" : 7,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-logging:jar",
    "to" : "org.apache.logging.log4j:log4j-to-slf4j:jar",
    "numericFrom" : 5,
    "numericTo" : 6,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-logging:jar",
    "to" : "org.slf4j:jul-to-slf4j:jar",
    "numericFrom" : 5,
    "numericTo" : 8,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter:jar",
    "to" : "org.springframework.boot:spring-boot-starter-logging:jar",
    "numericFrom" : 0,
    "numericTo" : 5,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter:jar",
    "to" : "jakarta.annotation:jakarta.annotation-api:jar",
    "numericFrom" : 0,
    "numericTo" : 9,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter:jar",
    "to" : "org.yaml:snakeyaml:jar",
    "numericFrom" : 0,
    "numericTo" : 10,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-actuator:jar",
    "to" : "org.springframework.boot:spring-boot-starter:jar",
    "numericFrom" : 11,
    "numericTo" : 0,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-actuator-autoconfigure:jar",
    "to" : "org.springframework.boot:spring-boot-actuator:jar",
    "numericFrom" : 12,
    "numericTo" : 13,
    "resolution" : "INCLUDED"
  }, {
    "from" : "com.fasterxml.jackson.core:jackson-databind:jar",
    "to" : "com.fasterxml.jackson.core:jackson-annotations:jar",
    "numericFrom" : 14,
    "numericTo" : 15,
    "resolution" : "INCLUDED"
  }, {
    "from" : "com.fasterxml.jackson.core:jackson-databind:jar",
    "to" : "com.fasterxml.jackson.core:jackson-core:jar",
    "numericFrom" : 14,
    "numericTo" : 16,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-actuator-autoconfigure:jar",
    "to" : "com.fasterxml.jackson.core:jackson-databind:jar",
    "numericFrom" : 12,
    "numericTo" : 14,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-actuator-autoconfigure:jar",
    "to" : "com.fasterxml.jackson.datatype:jackson-datatype-jsr310:jar",
    "numericFrom" : 12,
    "numericTo" : 17,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework:spring-context:jar",
    "to" : "org.springframework:spring-expression:jar",
    "numericFrom" : 18,
    "numericTo" : 19,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-actuator-autoconfigure:jar",
    "to" : "org.springframework:spring-context:jar",
    "numericFrom" : 12,
    "numericTo" : 18,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-actuator:jar",
    "to" : "org.springframework.boot:spring-boot-actuator-autoconfigure:jar",
    "numericFrom" : 11,
    "numericTo" : 12,
    "resolution" : "INCLUDED"
  }, {
    "from" : "io.micrometer:micrometer-core:jar",
    "to" : "org.hdrhistogram:HdrHistogram:jar",
    "numericFrom" : 20,
    "numericTo" : 21,
    "resolution" : "INCLUDED"
  }, {
    "from" : "io.micrometer:micrometer-core:jar",
    "to" : "org.latencyutils:LatencyUtils:jar",
    "numericFrom" : 20,
    "numericTo" : 22,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-actuator:jar",
    "to" : "io.micrometer:micrometer-core:jar",
    "numericFrom" : 11,
    "numericTo" : 20,
    "resolution" : "INCLUDED"
  }, {
    "from" : "com.macro.mall:mall:pom",
    "to" : "org.springframework.boot:spring-boot-starter-actuator:jar",
    "numericFrom" : 23,
    "numericTo" : 11,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework:spring-aop:jar",
    "to" : "org.springframework:spring-beans:jar",
    "numericFrom" : 24,
    "numericTo" : 25,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-aop:jar",
    "to" : "org.springframework:spring-aop:jar",
    "numericFrom" : 26,
    "numericTo" : 24,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-aop:jar",
    "to" : "org.aspectj:aspectjweaver:jar",
    "numericFrom" : 26,
    "numericTo" : 27,
    "resolution" : "INCLUDED"
  }, {
    "from" : "com.macro.mall:mall:pom",
    "to" : "org.springframework.boot:spring-boot-starter-aop:jar",
    "numericFrom" : 23,
    "numericTo" : 26,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "org.springframework.boot:spring-boot-test:jar",
    "numericFrom" : 28,
    "numericTo" : 29,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "org.springframework.boot:spring-boot-test-autoconfigure:jar",
    "numericFrom" : 28,
    "numericTo" : 30,
    "resolution" : "INCLUDED"
  }, {
    "from" : "net.minidev:accessors-smart:jar",
    "to" : "org.ow2.asm:asm:jar",
    "numericFrom" : 31,
    "numericTo" : 32,
    "resolution" : "INCLUDED"
  }, {
    "from" : "net.minidev:json-smart:jar",
    "to" : "net.minidev:accessors-smart:jar",
    "numericFrom" : 33,
    "numericTo" : 31,
    "resolution" : "INCLUDED"
  }, {
    "from" : "com.jayway.jsonpath:json-path:jar",
    "to" : "net.minidev:json-smart:jar",
    "numericFrom" : 34,
    "numericTo" : 33,
    "resolution" : "INCLUDED"
  }, {
    "from" : "com.jayway.jsonpath:json-path:jar",
    "to" : "org.slf4j:slf4j-api:jar",
    "numericFrom" : 34,
    "numericTo" : 35,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "com.jayway.jsonpath:json-path:jar",
    "numericFrom" : 28,
    "numericTo" : 34,
    "resolution" : "INCLUDED"
  }, {
    "from" : "jakarta.xml.bind:jakarta.xml.bind-api:jar",
    "to" : "jakarta.activation:jakarta.activation-api:jar",
    "numericFrom" : 36,
    "numericTo" : 37,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "jakarta.xml.bind:jakarta.xml.bind-api:jar",
    "numericFrom" : 28,
    "numericTo" : 36,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "org.assertj:assertj-core:jar",
    "numericFrom" : 28,
    "numericTo" : 38,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "org.hamcrest:hamcrest:jar",
    "numericFrom" : 28,
    "numericTo" : 39,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.junit.jupiter:junit-jupiter-api:jar",
    "to" : "org.opentest4j:opentest4j:jar",
    "numericFrom" : 40,
    "numericTo" : 41,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.junit.jupiter:junit-jupiter-api:jar",
    "to" : "org.junit.platform:junit-platform-commons:jar",
    "numericFrom" : 40,
    "numericTo" : 42,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.junit.jupiter:junit-jupiter:jar",
    "to" : "org.junit.jupiter:junit-jupiter-api:jar",
    "numericFrom" : 43,
    "numericTo" : 40,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.junit.jupiter:junit-jupiter:jar",
    "to" : "org.junit.jupiter:junit-jupiter-params:jar",
    "numericFrom" : 43,
    "numericTo" : 44,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.junit.jupiter:junit-jupiter:jar",
    "to" : "org.junit.jupiter:junit-jupiter-engine:jar",
    "numericFrom" : 43,
    "numericTo" : 45,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "org.junit.jupiter:junit-jupiter:jar",
    "numericFrom" : 28,
    "numericTo" : 43,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.junit.vintage:junit-vintage-engine:jar",
    "to" : "org.apiguardian:apiguardian-api:jar",
    "numericFrom" : 46,
    "numericTo" : 47,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.junit.vintage:junit-vintage-engine:jar",
    "to" : "org.junit.platform:junit-platform-engine:jar",
    "numericFrom" : 46,
    "numericTo" : 48,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.junit.vintage:junit-vintage-engine:jar",
    "to" : "junit:junit:jar",
    "numericFrom" : 46,
    "numericTo" : 49,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "org.junit.vintage:junit-vintage-engine:jar",
    "numericFrom" : 28,
    "numericTo" : 46,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.mockito:mockito-core:jar",
    "to" : "net.bytebuddy:byte-buddy:jar",
    "numericFrom" : 50,
    "numericTo" : 51,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.mockito:mockito-core:jar",
    "to" : "net.bytebuddy:byte-buddy-agent:jar",
    "numericFrom" : 50,
    "numericTo" : 52,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.mockito:mockito-core:jar",
    "to" : "org.objenesis:objenesis:jar",
    "numericFrom" : 50,
    "numericTo" : 53,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "org.mockito:mockito-core:jar",
    "numericFrom" : 28,
    "numericTo" : 50,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "org.mockito:mockito-junit-jupiter:jar",
    "numericFrom" : 28,
    "numericTo" : 54,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.skyscreamer:jsonassert:jar",
    "to" : "com.vaadin.external.google:android-json:jar",
    "numericFrom" : 55,
    "numericTo" : 56,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "org.skyscreamer:jsonassert:jar",
    "numericFrom" : 28,
    "numericTo" : 55,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework:spring-core:jar",
    "to" : "org.springframework:spring-jcl:jar",
    "numericFrom" : 57,
    "numericTo" : 58,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "org.springframework:spring-core:jar",
    "numericFrom" : 28,
    "numericTo" : 57,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "org.springframework:spring-test:jar",
    "numericFrom" : 28,
    "numericTo" : 59,
    "resolution" : "INCLUDED"
  }, {
    "from" : "org.springframework.boot:spring-boot-starter-test:jar",
    "to" : "org.xmlunit:xmlunit-core:jar",
    "numericFrom" : 28,
    "numericTo" : 60,
    "resolution" : "INCLUDED"
  }, {
    "from" : "com.macro.mall:mall:pom",
    "to" : "org.springframework.boot:spring-boot-starter-test:jar",
    "numericFrom" : 23,
    "numericTo" : 28,
    "resolution" : "INCLUDED"
  }, {
    "from" : "com.macro.mall:mall:pom",
    "to" : "cn.hutool:hutool-all:jar",
    "numericFrom" : 23,
    "numericTo" : 61,
    "resolution" : "INCLUDED"
  }, {
    "from" : "com.macro.mall:mall:pom",
    "to" : "org.projectlombok:lombok:jar",
    "numericFrom" : 23,
    "numericTo" : 62,
    "resolution" : "INCLUDED"
  }, {
    "from" : "com.macro.mall:mall:pom",
    "to" : "org.springframework.boot:spring-boot-configuration-processor:jar",
    "numericFrom" : 23,
    "numericTo" : 63,
    "resolution" : "INCLUDED"
  } ]
}
`
	f := must.A(os.CreateTemp("", ""))
	defer func() {
		must.Must(os.Remove(f.Name()))
	}()
	must.A(f.Write([]byte(a)))
	must.Must(f.Close())
	var d PluginGraphOutput
	must.Must(d.ReadFromFile(f.Name()))
	must.A(d.Tree())
}
