import { render, screen } from '@/utils/test-utils';
import { Layout } from './Layout';

describe('Layout', () => {
  it('renders children correctly', () => {
    render(
      <Layout>
        <div data-testid="content">Content</div>
      </Layout>
    );
    
    expect(screen.getByTestId('content')).toBeInTheDocument();
  });

  it('renders header when provided', () => {
    render(
      <Layout header={<div>Header</div>}>
        <div>Content</div>
      </Layout>
    );
    
    expect(screen.getByText('Header')).toBeInTheDocument();
  });
});