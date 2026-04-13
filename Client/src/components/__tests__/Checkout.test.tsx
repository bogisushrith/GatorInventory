import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { Checkout } from '../Checkout';
import { useCart } from '../../context/CartContext';
import { createOrder } from '../../api/orders';

const navigateMock = vi.fn();
const clearClientCartMock = vi.fn();
const refreshCartMock = vi.fn();

vi.mock('../../context/CartContext', () => ({
  useCart: vi.fn(),
}));

vi.mock('../../api/orders', () => ({
  createOrder: vi.fn(),
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

describe('Checkout', () => {
  beforeEach(() => {
    navigateMock.mockReset();
    clearClientCartMock.mockReset();
    refreshCartMock.mockReset();
    vi.mocked(createOrder).mockReset();

    vi.mocked(useCart).mockReturnValue({
      cartItems: [
        {
          id: 1,
          user_id: 1,
          product_id: 10,
          quantity: 2,
          product_name: 'Laptop',
          product_price: 1500,
          product_stock: 10,
        },
      ],
      loading: false,
      refreshCart: refreshCartMock,
      addItem: vi.fn(),
      updateItem: vi.fn(),
      removeItem: vi.fn(),
      clearClientCart: clearClientCartMock,
      totalQuantity: 2,
    });
  });

  it('displays checkout total', () => {
    render(<Checkout />);

    expect(screen.getByText('Laptop')).toBeInTheDocument();
    expect(screen.getByText('Total: $3000.00')).toBeInTheDocument();
  });

  it('handles place order click', async () => {
    vi.mocked(createOrder).mockResolvedValue({ order_id: 101 });

    render(<Checkout />);

    fireEvent.click(screen.getByRole('button', { name: /place order/i }));

    await waitFor(() => {
      expect(createOrder).toHaveBeenCalledWith([{ product_id: 10, quantity: 2 }]);
    });

    await waitFor(() => {
      expect(clearClientCartMock).toHaveBeenCalled();
      expect(refreshCartMock).toHaveBeenCalled();
      expect(navigateMock).toHaveBeenCalledWith('/history');
    });
  });
});
