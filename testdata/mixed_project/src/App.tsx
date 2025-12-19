import React from 'react';
import _ from 'lodash';

export const App = () => {
  const numbers = [1, 2, 3, 4, 5];
  const doubled = _.map(numbers, n => n * 2);
  return <div>Hello, {doubled}</div>;
};
