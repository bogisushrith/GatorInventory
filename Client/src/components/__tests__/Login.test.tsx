import { fireEvent, render, screen, waitFor } from '@testing-library/react';
import { Login } from '../Login';

const navigateMock = vi.fn();

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom');
  return {
    ...actual,
    useNavigate: () => navigateMock,
  };
});

describe('Login', () => {
  beforeEach(() => {
    navigateMock.mockReset();
    localStorage.clear();

    vi.stubGlobal(
      'fetch',
      vi.fn().mockResolvedValue({
        ok: true,
        json: async () => ({ role: 'user' }),
      })
    );
  });

  afterEach(() => {
    vi.unstubAllGlobals();
  });

  it('renders username and password input fields', () => {
    render(<Login />);

    expect(screen.getByLabelText(/username/i)).toBeInTheDocument();
    expect(screen.getByLabelText(/password/i)).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /login/i })).toBeInTheDocument();
  });

  it('accepts user input and submits login form', async () => {
    render(<Login />);

    fireEvent.change(screen.getByLabelText(/username/i), {
      target: { value: 'erkin' },
    });
    fireEvent.change(screen.getByLabelText(/password/i), {
      target: { value: 'Test1234' },
    });

    fireEvent.click(screen.getByRole('button', { name: /login/i }));

    await waitFor(() => {
      expect(fetch).toHaveBeenCalledWith(
        '/api/login',
        expect.objectContaining({
          method: 'POST',
          credentials: 'include',
        })
      );
    });

    await waitFor(() => {
      expect(navigateMock).toHaveBeenCalledWith('/dashboard');
    });
  });
});
