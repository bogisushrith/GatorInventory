import { fireEvent, render, screen } from '@testing-library/react';
import { ProductTable } from '../dashboard/ProductTable';
import { Product } from '../dashboard/types';

const mockProducts: Product[] = [
  { id: 1, name: 'Keyboard', price: 50, quantity: 2, category: 'Electronics' },
  { id: 2, name: 'Notebook', price: 10, quantity: 5, category: 'Stationery' }
];

describe('ProductTable', () => {
  // What this test does:
  // Verifies the table renders base heading and summary information from props.
  //
  // Why this is important:
  // Confirms users can see pagination summary and key table context at a glance.
  //
  // Steps:
  // 1. Render ProductTable with non-empty data.
  // 2. Assert the heading and "Showing X-Y of Z" text appear.
  it('renders heading and result summary', () => {
    render(
      <ProductTable
        products={mockProducts}
        loading={false}
        page={1}
        limit={5}
        totalResults={12}
        onLimitChange={vi.fn()}
        canManageProducts={true}
        isUser={false}
        cartQuantityByProductId={{}}
        cartDraftByProductId={{}}
        onCartDraftChange={vi.fn()}
        onAddToCart={vi.fn()}
        onApplyCartQuantity={vi.fn()}
        onUpdate={vi.fn()}
        onDelete={vi.fn()}
        stockDraftById={{}}
        onStockInputChange={vi.fn()}
        onUpdateStock={vi.fn()}
        updatingStockId={null}
      />
    );

    expect(screen.getByText('Products')).toBeInTheDocument();
    expect(screen.getByText('Showing 1–2 of 12 results')).toBeInTheDocument();
  });

  // What this test does:
  // Checks whether product values passed as props are shown correctly in the table.
  //
  // Why this is important:
  // Ensures backend data is actually visible to users and not lost in rendering.
  //
  // Steps:
  // 1. Render with mock products.
  // 2. Verify names, category values, and computed totals are visible.
  it('displays product rows from props', () => {
    render(
      <ProductTable
        products={mockProducts}
        loading={false}
        page={1}
        limit={5}
        totalResults={2}
        onLimitChange={vi.fn()}
        canManageProducts={true}
        isUser={false}
        cartQuantityByProductId={{}}
        cartDraftByProductId={{}}
        onCartDraftChange={vi.fn()}
        onAddToCart={vi.fn()}
        onApplyCartQuantity={vi.fn()}
        onUpdate={vi.fn()}
        onDelete={vi.fn()}
        stockDraftById={{}}
        onStockInputChange={vi.fn()}
        onUpdateStock={vi.fn()}
        updatingStockId={null}
      />
    );

    expect(screen.getByText('Keyboard')).toBeInTheDocument();
    expect(screen.getByText('Notebook')).toBeInTheDocument();
    expect(screen.getByText('Electronics')).toBeInTheDocument();
    expect(screen.getAllByText('Low Stock').length).toBeGreaterThan(0);
    expect(screen.getAllByText('$100.00')[0]).toBeInTheDocument();
  });

  // What this test does:
  // Ensures the empty-state message appears when there is no product data.
  //
  // Why this is important:
  // Prevents confusing blank UI and gives user feedback when filters return no items.
  //
  // Steps:
  // 1. Render with an empty products array.
  // 2. Confirm the empty-state message is shown.
  it('shows empty-state message when no products exist', () => {
    render(
      <ProductTable
        products={[]}
        loading={false}
        page={1}
        limit={5}
        totalResults={0}
        onLimitChange={vi.fn()}
        canManageProducts={false}
        isUser={true}
        cartQuantityByProductId={{}}
        cartDraftByProductId={{}}
        onCartDraftChange={vi.fn()}
        onAddToCart={vi.fn()}
        onApplyCartQuantity={vi.fn()}
        onUpdate={vi.fn()}
        onDelete={vi.fn()}
        stockDraftById={{}}
        onStockInputChange={vi.fn()}
        onUpdateStock={vi.fn()}
        updatingStockId={null}
      />
    );

    expect(screen.getByText('No products found. Try adjusting filters.')).toBeInTheDocument();
  });

  // What this test does:
  // Verifies interaction callbacks for edit/delete actions and items-per-page change.
  //
  // Why this is important:
  // Confirms UI controls trigger the right behavior hooks in parent logic.
  //
  // Steps:
  // 1. Render table with product data and admin actions enabled.
  // 2. Click Edit and Delete on first row.
  // 3. Change items-per-page selection.
  // 4. Assert corresponding handlers were called with expected payloads.
  it('triggers edit, delete and limit change callbacks', () => {
    const onUpdate = vi.fn();
    const onDelete = vi.fn();
    const onLimitChange = vi.fn();

    render(
      <ProductTable
        products={mockProducts}
        loading={false}
        page={1}
        limit={5}
        totalResults={2}
        onLimitChange={onLimitChange}
        canManageProducts={true}
        isUser={false}
        cartQuantityByProductId={{}}
        cartDraftByProductId={{}}
        onCartDraftChange={vi.fn()}
        onAddToCart={vi.fn()}
        onApplyCartQuantity={vi.fn()}
        onUpdate={onUpdate}
        onDelete={onDelete}
        stockDraftById={{ 1: '2', 2: '5' }}
        onStockInputChange={vi.fn()}
        onUpdateStock={vi.fn()}
        updatingStockId={null}
      />
    );

    fireEvent.click(screen.getAllByText('Edit')[0]);
    fireEvent.click(screen.getAllByText('Delete')[0]);
    fireEvent.change(screen.getAllByRole('combobox')[0], { target: { value: '10' } });

    expect(onUpdate).toHaveBeenCalledWith(mockProducts[0]);
    expect(onDelete).toHaveBeenCalledWith(1);
    expect(onLimitChange).toHaveBeenCalledWith(10);
  });
});
