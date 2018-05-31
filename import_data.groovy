#!/usr/bin/env groovy

@Grab('com.xlson.groovycsv:groovycsv:1.3')
import com.xlson.groovycsv.CsvParser

def file = new File(args[0])
def cvsData = new CsvParser().parse(file.text)

def records = cvsData.collect{ it ->
    def d = [:]
    d['name'] = it.Name
    d['email'] = it.Email
	d
}

println "DELETE FROM rsvp;"
println "DELETE FROM guests;"

def int invites = 0;
def int guests = 0;

records.findAll{ it.email }.groupBy{ it.email }.each{ k,v -> 
	def id = UUID.randomUUID().toString()
	println "INSERT INTO rsvp (rsvp_id, reply_type, reply_status, email) VALUES ('${id}', '', '', '${k}');"
	invites++
	v.each { g -> 
		println "INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('${id}',1,'${g.name}','');"
		guests++
	}
}

println "Invites: ${invites}, Guests: ${guests}"
