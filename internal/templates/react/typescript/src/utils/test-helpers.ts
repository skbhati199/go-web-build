import { act } from '@testing-library/react';

export const waitForAnimation = async () => {
  await act(async () => {
    await new Promise((resolve) => setTimeout(resolve, 0));
  });
};

export const mockIntersectionObserver = () => {
  const mockEntries = [{ isIntersecting: true }];
  const mockObserve = jest.fn();
  const mockUnobserve = jest.fn();
  const mockDisconnect = jest.fn();

  window.IntersectionObserver = jest.fn((callback) => ({
    observe: mockObserve,
    unobserve: mockUnobserve,
    disconnect: mockDisconnect,
    root: null,
    rootMargin: '',
    thresholds: [],
    takeRecords: () => mockEntries,
  }));

  return { mockObserve, mockUnobserve, mockDisconnect };
};