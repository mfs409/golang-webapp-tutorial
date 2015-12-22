// track calls to highlight that might interfere with other code blocks being
// loaded
var pendingHighlights = 0;

// highlight all of our code blocks
function highlightAllCode() {
    if (pendingHighlights !== 0)
        return;
    // go through the content, and for each code block, run the highlighter
    // on it if the highlighter hasn't been run before
    $('#content pre').each(function(i, pre) {
        var code = $(pre).find('code');
        if (!code.hasClass('hljs')) {
            hljs.highlightBlock(code.get(0));
        }
    });
}

// for each code snippet, request the file, parse it, escape any HTML special
// characters, and then dump the escaped code into the corresponding
// container
function grabSnippets() {
    $('code').each(function() {
        var block = $(this); // the code block
        var src = block.attr("source"); // get the href to the file to load into it
        if (!(src === undefined)) {
            pendingHighlights++;
            // make a data request
            $.get(src, function(data) {
                // escape the data, put it into the DOM, and call the highlighter
                var v = data.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;");
                block.text(data);
                pendingHighlights--;
                highlightAllCode();
            });
        }
    });
}

// this runs on document load, to highlight all code blocks and load all
// available styles.
$(document).ready(function() {
    grabSnippets();
});
