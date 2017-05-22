//lets globalize our json variables so that we on need one
//interval functions
var refreshrate = 500;
var systemObject = "";
var cpuObject = "";
var avg;
setInterval(function() {
  var xhr = new XMLHttpRequest();
  xhr.onload = function() {
    if (xhr.status == 200) {
      systemObject = JSON.parse(xhr.responseText);
      // console.log("UsedDiskPct:  "+systemObject.UsedDiskPct+"%");
    }
  };
  xhr.open('GET', 'https://localhost:8080/system');
  xhr.send(null);
  // console.log(systemObject)

  var Cpuxhr = new XMLHttpRequest();
  xhr.onload = function() {
    if (xhr.status == 200) {
      cpuObject = JSON.parse(Cpuxhr.responseText);
      // console.log(cpuObject.length);
      // for (i=0 ; i<cpuObject.length;i++){
      //   console.log("CPU " + i + ":  "+ cpuObject[i] + "%");
      // }
      let sum = cpuObject.reduce((previous, current) => current += previous);
      avg = sum / cpuObject.length;
      // console.log("total avg: " + avg)
    }
  };
  Cpuxhr.open('GET', 'https://localhost:8080/cpu');
  Cpuxhr.send(null);
  // console.log(cpuObject)


}, refreshrate);







///Include the needed packages
google.charts.load('current', {
  'packages': ['gauge']
});

///Set the function to draw the chart
google.charts.setOnLoadCallback(drawcpuChart);

///Function for the CPU chart
function drawcpuChart() {

  var cpudata = google.visualization.arrayToDataTable([
    ['Label', 'Value'],

    ['CPU', 55],

  ]);

  var options = {
    width: 150,
    height: 110,
    redFrom: 90,
    redTo: 100,
    yellowFrom: 75,
    yellowTo: 90,
    minorTicks: 20,
    animation: {
      duration: 50,
      easing: 'out',
    }
  };

  var cpuchart = new google.visualization.Gauge(document.getElementById('cpuchart_div'));

  cpuchart.draw(cpudata, options);

  setInterval(function() {
    cpudata.setValue(0, 1, avg);
    cpuchart.draw(cpudata, options);
  }, refreshrate);

}






//Function for the memmory Chart
google.charts.setOnLoadCallback(drawstoreageChart);

function drawstoreageChart() {

  var cpudata = google.visualization.arrayToDataTable([
    ['Label', 'Value'],
    ['Mem', 55],

  ]);

  var options = {
    width: 150,
    height: 110,
    redFrom: 90,
    redTo: 100,
    yellowFrom: 75,
    yellowTo: 90,
    minorTicks: 20,
    animation: {
    duration: 50,
    easing: 'out',
    }
  };

  var memChart = new google.visualization.Gauge(document.getElementById('memchart_div'));

  memChart.draw(cpudata, options);

  setInterval(function() {
    var xhr = new XMLHttpRequest();
    xhr.onload = function() {
      if (xhr.status == 200) {

        systemObject = JSON.parse(xhr.responseText);
        // console.log(percent+"%");

      }
    };
    xhr.open('GET', 'https://localhost:8080/system');
    xhr.send(null);
    // console.log(systemObject.Os)

    cpudata.setValue(0, 1, systemObject.MemUsed);
    memChart.draw(cpudata, options);
  }, refreshrate);
}


google.charts.setOnLoadCallback(drawHDChart);

//Function for the memmory Chart
function drawHDChart() {

  var cpudata = google.visualization.arrayToDataTable([
    ['Label', 'Value'],
    ['HD', 55],

  ]);

  var options = {
    width: 150,
    height: 420,
    redFrom: 90,
    redTo: 100,
    yellowFrom: 75,
    yellowTo: 90,
    minorTicks: 20,
    animation: {
    duration: 50,
    easing: 'out',
    }
  };
  ///set the document object
  var cpuchart = new google.visualization.Gauge(document.getElementById('storeagechart_div'));

  cpuchart.draw(cpudata, options);

  setInterval(function() {
    cpudata.setValue(0, 1, systemObject.UsedDiskPct);
    cpuchart.draw(cpudata, options);
  }, refreshrate);
}


///set the collapsible menus to 'open'
$('.collapsible').collapsible('open', 0);



//initialize some global variables
var label = [""];
var seriesData = [
  [5],
  [3],
  [15],
  [9]

]

var chartdata = {
  // A labels array that can contain any sort of values
  labels: label,
  // Our series array that contains series objects or in this case series data arrays
  series: seriesData,
  high: 100,
  low: 0,
};

// Create a new line chart object where as first parameter we pass in a selector
// that is resolving to our chart container element. The Second parameter
// is the actual data object.

setInterval(function() {
  new Chartist.Line('.ct-chart', chartdata);
  var d = new Date();
  var time = d.toISOString();

  label.push(time);
  if (label.length > 10) {
    label.shift();
  }

  //find the number of graphlines
  for (j = 0; j <= seriesData.length; j++) {
    //find the number of elements in the array
    for (i = 0; i <= seriesData[j].length; i++) {
      seriesData[i].push(parseInt(cpuObject[i]));
      console.log("cpu" + i + ": " + parseInt(cpuObject[i]));
      if (seriesData[i].length > 10) {
        seriesData[i].shift();

      }
    }
  }


  for (i = 0; i < seriesData.length; i++) {
    if (seriesdata[i].length > 10) {
      seriesdata[i].shift();

    }
  }

}, refreshrate);
