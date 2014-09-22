
INSTRUCTIONS
————————————

$ <— connotes terminal start line


1) In the terminal, navigate to the directory of this file:
   $ cd /path/to/downloaded/simple_web_app/


2) Build the app from root directory of this app with go (make sure go is installed):
   $ go build simple_web_app.go


3) Run the newly created build from the same directory:
   $ ./simple_web_app


4) In any web browser, enter the following URL:
   http://localhost:8080/select

5) Select the graphs you want to see, then click “View Graphs” button.  You’re off!


OPTIONS
—————————

Ids: network, memory, temperature, disk, cpu, energy


Get Graph Data (URL):
  http://localhost:8080/data/?type=[graphId1]&[graphId2]

  EXAMPLE: http://localhost:8080/data/?type=network&cpu


View Graphs Option 1 (GUI):
  —Go to the following url in any web browser: http://localhost:8080/select

View Graphs Option 2 (URL):
  http://localhost:8080/display/?type=[graphId1]&[graphId2]

  EXAMPLE: http://localhost:8080/display/?type=network&cpu





