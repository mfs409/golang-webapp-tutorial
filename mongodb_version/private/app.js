// All our templates go in here
var hb_t = {};

// This list holds all of the template names.  Names matter... for entry
// "abc" in this list, we will create field hb_t.abc in the above map, using
// file "private/abc.hb"
var hblist = [ "dataList", "editModal", "addModal", "nav" ];

// When the page has been loaded, and we know all the javascript code is
// available, call loadTemplates to fetch and parse the templates.  When all
// templates are loaded, buildPage will run
$(document).ready(function() {
    loadTemplates(buildPage);
});

// Load the handlebar template files.  Once all the templates are loaded, run "action"
function loadTemplates(action) {
    // keep track of the number of templates we still need to load
    var remain = hblist.length;
    // Function to load template then call action() iff all templates loaded
    var loader = function(i) {
        $.get("private/" + hblist[i]+".hb", function(data) {
            hb_t[hblist[i]] = Handlebars.compile(data);
            remain--;
            if (remain == 0) action();
        }, "html");
    }
    // load the templates
    for (var i = 0; i < hblist.length; i++) { loader(i); }
}

// We build the main logged-in page by adding the initial components, which
// are a navbar and the data list
function buildPage() {
    navbar.Inject();
    dataList.DisplayLatest();
}
