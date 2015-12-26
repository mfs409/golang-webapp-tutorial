<!DOCTYPE html>
<html lang="en">
 <head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1">

  <meta name="description"
        content="Example Golang Web App with MySQL and OAuth 2.0">
  <meta name="author" content="mfs409@gmail.com">

  <title>Example Web App</title>

  <link href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/css/bootstrap.min.css" rel="stylesheet" integrity="sha256-7s5uDGW3AHqw6xtJmNNtr+OBRJUlgkNJEo78P4b0yRw= sha512-nNo+yCHEyn0smMxSswnf/OnX6/KwJuZTlNZBjauKhTK0c+zT+q5JOCx0UFhXQ6rJR9jg6Es8gPuD2uZcYDLqSw==" crossorigin="anonymous">
  <script src="https://code.jquery.com/jquery-2.1.4.min.js"></script>
  <script src="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.6/js/bootstrap.min.js" integrity="sha256-KXn5puMvxCw+dAYznun+drMdG1IFl3agK0p/pqT9KAo= sha512-2e8qq0ETcfWRI4HJBzQiA3UoyFk6tbNyG+qSaIBZLyW9Xf3sWZHN/lxe9fTh1U45DpPf07yj94KsUHHWe4Yk1A==" crossorigin="anonymous"></script>
  <script src="/private/app.js"></script>
 </head>

 <body role="document">

  <div class="container" role="main">

   <div class="jumbotron">
    <h1>Example Web App</h1>
       <a class="btn btn-default" href="/logout">
        <span class="glyphicon glyphicon-log-out"> Log-Out</span>
       </a>
   </div>

   <div class="row">
    <div class="col-lg-12">
     <table class="table table-striped" id="dataTable">
      <thead>
       <tr>
        <th>Small Note</th>
        <th>Big Note</th>
        <th>Favorite Int</th>
        <th>Favorite Float</th>
        <th>Extra Float?</th>
        <th><a href='#editModal' data-toggle='modal' title='Add'
               class='btn btn-default btn-sm add-btn' data-id=''><span class='glyphicon glyphicon-plus'></span></a></th>
       </tr>
      </thead>
      <tbody>
      </tbody>
     </table>
    </div>
   </div>
  </div>

  <div id="editModal" class="modal fade" role="dialog">
   <div class="modal-dialog">
    <div class="modal-content">
     <div class="modal-header">
      <button type="button" class="close" data-dismiss="modal">
       <span class="glyphicon glyphicon-remove"></span>
      </button>
      <h4 id="modal-header" class="modal-title">Edit Entry</h4>
     </div>
     <div class="modal-body">
      <input type="hidden" class="" id="id" placeholder="-">
      <div class="form-group">
       <label for="smallnote">Small Note</label>
       <input type="text" class="form-control" id="smallnote" placeholder="-">
      </div>
      <div class="form-group">
       <label for="bignote">Big Note</label>
       <input type="text" class="form-control" id="bignote" placeholder="-">
      </div>
      <div class="form-group">
       <label for="favint">Favorite Integer</label>
       <input type="text" class="form-control" id="favint" placeholder="-">
      </div>
      <div class="form-group">
       <label for="favfloat">Favorite Float</label>
       <input type="text" class="form-control" id="favfloat" placeholder="-">
      </div>
      <div class="form-group">
       <label for="trickfloat">Tricky Float</label>
       <input type="text" class="form-control" id="trickfloat" placeholder="-">
      </div>
     </div>
     <div class="modal-footer">
      <button id="modal-save" type="button" class="btn btn-success" data-dismiss="modal">Save</button>
      <button type="button" class="btn btn-default" data-dismiss="modal">Close</button>
     </div>
    </div>
   </div>
  </div>
 </body>
</html>
