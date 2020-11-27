the micro service listens to http://0.0.0.0:8000/
it takes json file by post method at url "/"
in result return json containing total flight path. 
example input:
{
    "flights":[
       [
          "IND",
          "EWR"
       ],
       [
          "SFO",
          "ATL"
       ],
       [
          "GSO",
          "IND"
       ],
       [
          "ATL",
          "GSO"
       ]
    ]
 }

example output:
{"FlightPath":["SFO","EWR"]}

# run program by 
go run main.go


# to test with curl, run command from prject directory:
curl  -X POST http://127.0.0.1:8000/ -d @test1.json -H "Content-Type: application/json"
curl  -X POST http://127.0.0.1:8000/ -d @test2.json -H "Content-Type: application/json"
curl  -X POST http://127.0.0.1:8000/ -d @test3.json -H "Content-Type: application/json"


# to buils docker image run:
docker build -t flightpath .

# run docker container
docker run -i -t -p 8000:8000 flightpath