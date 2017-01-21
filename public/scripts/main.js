/**
 * Created by kanta on 1/21/17.
 */
document.addEventListener("DOMContentLoaded", function () {
    var data = [
            ["", "Ford", "Volvo", "Toyota", "Honda"],
            ["2016", 10, 11, 12, 13],
            ["2017", 20, 11, 14, 13],
            ["2018", 30, 15, 12, 13]
        ],
        container = document.getElementById("example"),
        searchField = document.getElementById('search_field'),
        resultCount = document.getElementById('resultCount'),
        hot,
        searchResultCount = 0;

    //Search
    function searchResultCounter(instance, row, col, value, result) {
        Handsontable.Search.DEFAULT_CALLBACK.apply(this, arguments);

        if (result) {
            searchResultCount++;
        }
    }

    hot = new Handsontable(container, {
        afterChange: function (changes, source) {
            console.log("This is fired after change.");
            console.log("CHANGES: ", changes);
            console.log("SOURCE: ", source);
        },
        //afterChange se triggera kada se dogodi promjena.
        afterCreateRow: function (index, amount, source) {
            console.log("index: ", index);
            console.log("amount: ", amount);
            console.log("SOURCE: ", source);
        },
        search: {
            callback: searchResultCounter
        },
        data: data,
        minSpareCols: 1,
        minSpareRows: 1,
        rowHeaders: true,
        colHeaders: true,
        contextMenu: true,
        //bla
        manualRowResize: true,
        columnSorting: true,
        manualColumnResize: true,
        maxRows: 7,
        maxCols: 7
    });

    Handsontable.Dom.addEvent(searchField, 'keyup', function(event) {
        var queryResult;

        searchResultCount = 0;
        queryResult = hot.search.query(this.value);
        console.log("queryResult", queryResult);
        resultCount.innerText = searchResultCount.toString();
        hot.render();
    });

    // This will be usefull for getting data dump from the whole table.
    /*
    function bindDumpButton() {
        if (typeof Handsontable === "undefined") {
            return;
        }

        Handsontable.Dom.addEvent(document.body, 'click', function (e) {

            var element = e.target || e.srcElement;

            if (element.nodeName == "BUTTON" && element.name == 'dump') {
                var name = element.getAttribute('data-dump');
                var instance = element.getAttribute('data-instance');
                var hot = window[instance];
                console.log('data of ' + name, hot.getData());
            }
        });
    }
    bindDumpButton();
    */
});
