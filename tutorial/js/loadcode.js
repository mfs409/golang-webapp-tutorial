// this runs on document load, to fetch and highlight all code blocks
$(document).ready(function() {
    // for each code snippet, fetch the URL (via 'source' sttribute), escape
    // any HTML special characters, and then dump the escaped code into the
    // corresponding container
    $('pre code').each(function(i, block) {
        var src = $(block).attr('source');
        if (src !== undefined) 
            $.get(src, function(data) {
                $(block).text(data.replace(/&/g, "&amp;").replace(/</g, "&lt;").replace(/>/g, "&gt;"));
                hljs.highlightBlock(block);
            });
        else 
            hljs.highlightBlock(block);
    });
});
