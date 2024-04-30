console.log(js_data);
// set the dimensions and margins of the graph
var width = 450, height = 450, margin = 40;
// The radius of the pieplot is half the width or half the height (smallest one). I subtract a bit of margin.
var radius = Math.min(width, height) / 2 - margin;
// append the svg object to the div called 'my_dataviz'
var svg = d3.select("#my_dataviz")
    .append("svg")
    .attr("width", width)
    .attr("height", height)
    .append("g")
    .attr("transform", "translate(".concat(width / 2, ",").concat(height / 2, ")"));
// Create dummy data
var ram_data = { a: js_data.RAM, b: js_data.RAMInUse };
// set the color scale
var color = d3.scaleOrdinal()
    .range(["#ffffff", "#3b75cb"]);
// Compute the position of each group on the pie:
var pie = d3.pie()
    .value(function (d) { return d[1]; });
var data_ready = pie(Object.entries(ram_data));
// Build the pie chart: Basically, each part of the pie is a path that we build using the arc function.
svg
    .selectAll('whatever')
    .data(data_ready)
    .join('path')
    .attr('d', d3.arc()
    .innerRadius(100) // This is the size of the donut hole
    .outerRadius(radius))
    .attr('fill', function (d) { return color(d.data[0]); })
    .attr("stroke", "black")
    .style("stroke-width", "2px")
    .style("opacity", 0.7);
