// The nav object is a singleton web component for providing a simple
// bootstrap navbar on the page.
//
// NB: This doesn't really need to be a component with an Handlebars
//     template, but it's easier to work this way.  If we embedded the nav in
//     the app.tpl file, we'd have to re-start the web server any time we
//     edited the navbar.  Also, it's easier to keep track of what's been put
//     in the DOM when the document body is pretty much empty to start with.
var navbar = {

    // Put the navbar into the DOM as the first child of the body
    Inject : function() {
        $("body").prepend(hb_t.nav());
    }
    
}
