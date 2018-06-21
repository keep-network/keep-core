import React, { Component } from 'react'
import { Line } from 'react-chartjs-2'
import moment from 'moment'

class VestingChart extends Component {
  render() {
    if (this.props.details) {
      return <Line data={ this.chartData() } options={ this.chartOptions() } />
    } else {
      return "Loading ...";
    }
  }

  chartData() {
    return {
      datasets: [
        this.fromBaseDataset({
          data: this.getPoints()
        }),
      ],
    }
  }

  getPoints() {
    const { start, cliff, end } = this.props.details

    const now = new Date() / 1000 // normalize to seconds

    const points = [ this.getDataPointAt(start) ]

    if (cliff < now) {
      points.push(this.getDataPointAt(cliff))
    }

    if (start < now && now < end) {
      points.push(this.getDataPointAt(now))
    }

    if (cliff > now) {
      points.push(this.getDataPointAt(cliff))
    }

    points.push(this.getDataPointAt(end))

    return points
  }

  getDataPointAt(date) {
    return {
      x: this.formatDate(date),
      y: this.getAmountAt(date)
    }
  }

  formatDate(date) {
    return moment(date * 1000).format('MM/DD/YYYY HH:mm')
  }

  getAmountAt(date) {
    const { amount, start, end, decimals } = this.props.details
    const slope = (date - start) / (end - start)
    return this.displayAmount(amount, decimals) * slope
  }

  displayAmount(amount, decimals) {
    amount = amount / (10 ** decimals)
    return Math.round(amount * 10000) / 10000
  }

  chartOptions() {
    return {
      legend: { display: false },
      scales: {
        xAxes: [ {
          type: "time",
          time: {
            format: 'MM/DD/YYYY HH:mm',
            tooltipFormat: 'll HH:mm'
          },
          scaleLabel: {
            display: true,
            labelString: 'Date'
          }
        }, ],
        yAxes: [{
          scaleLabel: {
            display: true,
            labelString: this.props.details.symbol || ''
          }
        }]
      },
    }
  }

  fromBaseDataset(opts) {
    return {
      lineTension: 0.1,
      backgroundColor: 'rgba(92,182,228,0.4)',
      borderColor: 'rgba(92,182,228,1)',
      borderJoinStyle: 'miter',
      pointBorderColor: 'rgba(92,182,228,1)',
      pointBackgroundColor: 'rgba(92,182,228,1)',
      pointBorderWidth: 1,
      pointHoverRadius: 5,
      pointHoverBackgroundColor: 'rgba(92,182,228,1)',
      pointHoverBorderColor: 'rgba(220,220,220,1)',
      pointHoverBorderWidth: 2,
      pointRadius: 5,
      pointHitRadius: 10,
      ...opts
    }
  }
}

export default VestingChart
