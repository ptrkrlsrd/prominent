# prominent
## Microservice for getting the n most prominent colors of an image

## Routes
* POST - /api/v1/analyze

## Usage
* Start the program
* Post an image to the API
  * Example: `curl -F "image=@~/Images/image.jpg" localhost:3000/api/v1/analyze&n=5`
