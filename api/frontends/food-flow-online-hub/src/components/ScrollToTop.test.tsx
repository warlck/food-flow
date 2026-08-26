import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/react';
import React from 'react';
import { ScrollToTop } from './ScrollToTop';
import * as reactRouter from 'react-router-dom';

vi.mock('react-router-dom', () => ({
  useLocation: vi.fn(),
  useNavigationType: vi.fn(),
  NavigationType: {
    Pop: 'POP',
    Push: 'PUSH',
    Replace: 'REPLACE',
  },
}));

describe('ScrollToTop', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    window.scrollTo = vi.fn();
  });

  it('scrolls to top on PUSH navigation without hash', () => {
    vi.mocked(reactRouter.useLocation).mockReturnValue({
      pathname: '/contact',
      hash: '',
      search: '',
      state: null,
      key: 'test1',
    });
    vi.mocked(reactRouter.useNavigationType).mockReturnValue(reactRouter.NavigationType.Push as unknown as ReturnType<typeof reactRouter.useNavigationType>);

    render(<ScrollToTop />);

    expect(window.scrollTo).toHaveBeenCalledTimes(1);
    expect(window.scrollTo).toHaveBeenCalledWith({ top: 0, behavior: 'auto' });
  });

  it('scrolls to top on REPLACE navigation without hash', () => {
    vi.mocked(reactRouter.useLocation).mockReturnValue({
      pathname: '/privacy',
      hash: '',
      search: '',
      state: null,
      key: 'test2',
    });
    vi.mocked(reactRouter.useNavigationType).mockReturnValue(reactRouter.NavigationType.Replace as unknown as ReturnType<typeof reactRouter.useNavigationType>);

    render(<ScrollToTop />);

    expect(window.scrollTo).toHaveBeenCalledTimes(1);
    expect(window.scrollTo).toHaveBeenCalledWith({ top: 0, behavior: 'auto' });
  });

  it('does NOT scroll on POP navigation (back / forward)', () => {
    vi.mocked(reactRouter.useLocation).mockReturnValue({
      pathname: '/contact',
      hash: '',
      search: '',
      state: null,
      key: 'test3',
    });
    vi.mocked(reactRouter.useNavigationType).mockReturnValue(reactRouter.NavigationType.Pop as unknown as ReturnType<typeof reactRouter.useNavigationType>);

    render(<ScrollToTop />);

    expect(window.scrollTo).not.toHaveBeenCalled();
  });

  it('does NOT scroll when URL has a hash anchor', () => {
    vi.mocked(reactRouter.useLocation).mockReturnValue({
      pathname: '/',
      hash: '#merchant-cta',
      search: '',
      state: null,
      key: 'test4',
    });
    vi.mocked(reactRouter.useNavigationType).mockReturnValue(reactRouter.NavigationType.Push as unknown as ReturnType<typeof reactRouter.useNavigationType>);

    render(<ScrollToTop />);

    expect(window.scrollTo).not.toHaveBeenCalled();
  });
});
