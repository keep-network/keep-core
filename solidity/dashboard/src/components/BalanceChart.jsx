import React, { Component } from 'react';
import { Pie } from 'react-chartjs-2';
import moment from 'moment';

class BalanceChart extends Component {
  render() {
    return <Pie data={ this.chartData() } options={ this.chartOptions() } />;
  }

  chartData() {
    return {
      datasets: [
        this.fromBaseDataset({
          data: [
          ]
        }),
      ],
    }
  }

  formatDate(date) {
    return moment(date * 1000).format('MM/DD/YYYY HH:mm');
  }

  chartOptions() {
    return {
      legend: { display: true },
    }
  }

  fromBaseDataset(opts) {
    return {
    }
  }
}

export default BalanceChart