import { render, screen, waitFor } from '@testing-library/react';
import { OrderHistory } from '../OrderHistory.tsx';
import { getOrders } from '../../api/orders';

const navigateMock = vi.fn();

vi.mock('../../api/orders', () => ({
  getOrders: vi.fn(),
}));

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom');
  return {
    ...actual,
    useNavigate: () => navigateMock,
  };
});

vi.mock('../dashboard/DashboardShell', () => ({
  DashboardShell: ({ children }: { children: unknown }) => <div>{children as any}</div>,
}));

describe('OrderHistory', () => {
  beforeEach(() => {
    navigateMock.mockReset();
    vi.mocked(getOrders).mockReset();
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('displays orders list', async () => {
    vi.mocked(getOrders).mockResolvedValue([
      {
        id: 101,
        user_id: 1,
        user_name: 'erkin',
        status: 'pending',
        created_at: '2026-04-12T10:00:00Z',
        items: [{ id: 1, order_id: 101, product_id: 10, product_name: 'Laptop', product_price: 1500, quantity: 2 }],
      },
    ]);

    render(<OrderHistory />);

    await waitFor(() => {
      expect(screen.getByText('Order #101')).toBeInTheDocument();
      expect(screen.getByText('Laptop')).toBeInTheDocument();
    });
  });

  it('shows empty state when there are no orders', async () => {
    vi.mocked(getOrders).mockResolvedValue([]);

    render(<OrderHistory />);

    await waitFor(() => {
      expect(screen.getByText('No orders yet')).toBeInTheDocument();
    });
  });
});
