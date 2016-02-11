{{! This file provides a template for generating the navbar}}
{{! It's just a Bootstrap navbar... see the Bootstrap docs for ideas}}
{{! about how to make it fancier}}
<nav class="navbar navbar-default">
 <div class="container-fluid">
  
  <div class="navbar-header">
   <button type="button" class="navbar-toggle collapsed"
           data-toggle="collapse" data-target="#nav-collapsable"
           aria-expanded="false">
    <span class="sr-only">Toggle navigation</span>
    <span class="icon-bar"></span>
    <span class="icon-bar"></span>
    <span class="icon-bar"></span>
   </button>
   <a class="navbar-brand" href="#">Example Web App</a>
  </div>

  <div class="collapse navbar-collapse" id="nav-collapsable">
   <ul class="nav navbar-nav navbar-right">
    <li>
     <a href="/logout">
      <span class="glyphicon glyphicon-log-out"></span> Log-Out
     </a>
    </li>
   </ul>
  </div>
  
 </div>
</nav>
