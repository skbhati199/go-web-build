const { performance, PerformanceObserver } = require('perf_hooks');
const webpack = require('webpack');

class BuildPerformanceMonitor {
  constructor() {
    this.observer = new PerformanceObserver((list) => {
      const entries = list.getEntries();
      entries.forEach((entry) => {
        console.log(`${entry.name}: ${entry.duration}ms`);
      });
    });
    
    this.observer.observe({ entryTypes: ['measure'], buffered: true });
  }

  startMeasure(phase) {
    performance.mark(`${phase}-start`);
  }

  endMeasure(phase) {
    performance.mark(`${phase}-end`);
    performance.measure(phase, `${phase}-start`, `${phase}-end`);
  }
}

module.exports = BuildPerformanceMonitor;