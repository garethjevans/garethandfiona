#!/usr/bin/env groovy

@Grab('com.xlson.groovycsv:groovycsv:1.3')
import com.xlson.groovycsv.CsvParser

def file = new File("/scripts/guests.csv")
def cvsData = new CsvParser().parse(file.text)

def records = cvsData.collect{ it ->
	// Name,email address,day or evening,unlikely,Save the date,Invite sent,RSVP,accomodation booked,,,,,Maybes
    def d = [:]
    d['name'] = it.Name
    d['email'] = it.'email address'
    d['unlikely'] = (it.'unlikely' == 'y')
	d
}

println "DELETE FROM rsvp;"
println "DELETE FROM guests;"

def int invites = 0;
def int guests = 0;

records.findAll{ !it.unlikely }.groupBy{ it.email }.each{ k,v -> 
	def id = UUID.randomUUID().toString()
	println "INSERT INTO rsvp (rsvp_id, email) VALUES ('${id}', '${k}');"
	invites++
	v.each { g -> 
		println "INSERT INTO guests (rsvp_id, attending, name, comments) VALUES ('${id}',1,'${g.name}','');"
		guests++
	}
}

println "Invites: ${invites}, Guests: ${guests}"
