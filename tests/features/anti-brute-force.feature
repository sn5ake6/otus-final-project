# file: features/anti-brute-force.feature

# http://anti-brute-force:8888/

Feature: anti-brute-force
	In order to test application behavior
	As API client of service
	I need to be able to run features

	Scenario: Authorize success
		When I send "POST" request to "http://anti-brute-force:8888/authorize/" with "application/json" data:
		"""
		{
			"login":"login",
			"password":"password",
			"ip":"127.0.0.1"
		}
		"""
		Then The response code should be 200
		And The response should match text:
		"""
		{
			"ok": true
		}
		"""

	Scenario: Add to blacklist
		When I send "POST" request to "http://anti-brute-force:8888/blacklist/" with "application/json" data:
		"""
		{
			"subnet": "192.1.1.0/25"
		}
		"""
		Then The response code should be 201

	Scenario: Authorize from blacklist
		When I send "POST" request to "http://anti-brute-force:8888/authorize/" with "application/json" data:
		"""
		{
			"login":"blacklist login",
			"password":"blacklist password",
			"ip":"192.1.1.2"
		}
		"""
		Then The response code should be 429
		And The response should match text:
		"""
		{
			"ok": false
		}
		"""

	Scenario: Delete from blacklist
		When I send "DELETE" request to "http://anti-brute-force:8888/blacklist/" with "application/json" data:
		"""
		{
			"subnet": "192.1.1.0/25"
		}
		"""
		Then The response code should be 200

	Scenario: Add to whitelist
		When I send "POST" request to "http://anti-brute-force:8888/whitelist/" with "application/json" data:
		"""
		{
			"subnet": "192.1.1.0/25"
		}
		"""
		Then The response code should be 201


	Scenario: Authorize from whitelist
		When I send "POST" request to "http://anti-brute-force:8888/authorize/" with "application/json" data:
		"""
		{
			"login":"whitelist login",
			"password":"whitelist password",
			"ip":"192.1.1.2"
		}
		"""
		Then The response code should be 200
		And The response should match text:
		"""
		{
			"ok": true
		}
		"""

	Scenario: Delete from whitelist
		When I send "DELETE" request to "http://anti-brute-force:8888/whitelist/" with "application/json" data:
		"""
		{
			"subnet": "192.1.1.0/25"
		}
		"""
		Then The response code should be 200

	Scenario: Repeated authorize fail
		When I send "POST" request to "http://anti-brute-force:8888/authorize/" and repeat it 10 times with "application/json" data:
		"""
		{
			"login":"login",
			"password":"password",
			"ip":"192.1.1.2"
		}
		"""
		Then The response code should be 429
		And The response should match text:
		"""
		{
			"ok": false
		}
		"""

	Scenario: Authorize success after wait reset interval
		When I wait "1m" and send "POST" request to "http://anti-brute-force:8888/authorize/" with "application/json" data:
		"""
		{
			"login":"login",
			"password":"password",
			"ip":"192.1.1.2"
		}
		"""
		Then The response code should be 200
		And The response should match text:
		"""
		{
			"ok": true
		}
		"""