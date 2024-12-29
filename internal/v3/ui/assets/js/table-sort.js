function sortTable(tableId, columnIndex) {
    var table = document.getElementById(tableId);
    var switching = true;
    var dir = "asc";
    var switchcount = 0;
    
    while (switching) {
        switching = false;
        var rows = table.rows;
        
        for (var i = 1; i < (rows.length - 1); i++) {
            var shouldSwitch = false;
            var x = rows[i].getElementsByTagName("TD")[columnIndex];
            var y = rows[i + 1].getElementsByTagName("TD")[columnIndex];
            
            var xContent = x.innerHTML.toLowerCase();
            var yContent = y.innerHTML.toLowerCase();
            
            // Try to convert to numbers if the content looks numeric
            if (xContent.match(/^[\d./]+$/) && yContent.match(/^[\d./]+$/)) {
                xContent = parseFloat(xContent.replace(/[^\d.-]/g, '')) || xContent;
                yContent = parseFloat(yContent.replace(/[^\d.-]/g, '')) || yContent;
            }
            
            if (dir === "asc") {
                if (xContent > yContent) {
                    shouldSwitch = true;
                    break;
                }
            } else if (dir === "desc") {
                if (xContent < yContent) {
                    shouldSwitch = true;
                    break;
                }
            }
        }
        
        if (shouldSwitch) {
            rows[i].parentNode.insertBefore(rows[i + 1], rows[i]);
            switching = true;
            switchcount++;
        } else {
            if (switchcount === 0 && dir === "asc") {
                dir = "desc";
                switching = true;
            }
        }
    }
} 