<!DOCTYPE html>
<html lang="en">
 <head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">

  <meta name="description"
        content="Example Golang Web App with MongoDB and OAuth 2.0">
  <meta name="author" content="mfs409@gmail.com">

  <title>Example Web App</title>

  <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" rel="stylesheet">
  <script src="https://code.jquery.com/jquery-2.1.4.min.js"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js"></script>
 </head>

 <body role="document">

  <div class="container" role="main">

   <div class="jumbotron">
    <h1>Example Web App</h1>
   </div>

   {{if .Inf}}
   <div class="alert alert-info alert-dismissible" role="alert">
    <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
    {{.InfText}}
   </div>
   {{end}}
   
   {{if .Err}}
   <div class="alert alert-danger alert-dismissible" role="alert">
    <button type="button" class="close" data-dismiss="alert" aria-label="Close"><span aria-hidden="true">&times;</span></button>
    {{.ErrText}}
   </div>
   {{end}}
   
   <div class="row">
    <div class="col-sm-6">
     <div class="panel panel-default">
      <div class="panel-heading">
       <h3 class="panel-title">Register</h3>
      </div>
      <div class="panel-body">
       <p>
        This site uses Google IDs for authentication.  Use the following link
        to register your ID with our system.
       </p>
       <a class="btn btn-default" href="/register">
        <span class="glyphicon glyphicon-user"></span> Register
       </a>
      </div>
     </div>
    </div>
    <div class="col-sm-6">
     <div class="panel panel-default">
      <div class="panel-heading">
       <h3 class="panel-title"> Log In</h3>
      </div>
      <div class="panel-body">
       <p>
        If you have already registered, you can use this link to log in,
        using your Google ID.
       </p>
       <a class="btn btn-default" href="/login">
        <span class="glyphicon glyphicon-log-in"></span> Log-In
       </a>
      </div>
     </div>
    </div>
   </div>
  </div>
 </body>
</html>
