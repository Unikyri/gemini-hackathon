import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render } from '@testing-library/react';
import * as fc from 'fast-check';
import { NodeWorkspace } from '../components/NodeWorkspace';
import * as hooks from '@/shared/hooks';
import type { GetNodeResponse } from '@/shared/api';

// Mock the useNodeDetail hook
vi.mock('@/shared/hooks', async () => {
  const actual = await vi.importActual('@/shared/hooks');
  return {
    ...actual,
    useNodeDetail: vi.fn(),
  };
});

/**
 * Property 6: Node Title Display Consistency
 *
 * For any successfully loaded node, the title displayed in the Workspace header
 * should be exactly equal to the title field from the API response.
 */
describe('NodeWorkspace - Property-Based Tests', () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe('Property 6: Node Title Display Consistency', () => {
    it('should display the exact title from the API response for any valid title string', () => {
      fc.assert(
        fc.property(
          // Generate arbitrary title strings
          fc.oneof(
            fc.string({ minLength: 1, maxLength: 200 }),
            fc.string({ minLength: 1, maxLength: 100 }).map(s => `Exercise: ${s}`),
            fc.string({ minLength: 1, maxLength: 100 }).map(s => `${s} - Part 1`),
            fc.constantFrom(
              'Simple Title',
              'Title with numbers 123',
              'Title with special chars !@#$%',
              'Very Long Title That Spans Multiple Words And Contains Lots Of Information About The Exercise',
              'Title\nwith\nnewlines',
              'Title\twith\ttabs'
            )
          ),
          fc.string({ minLength: 10, maxLength: 500 }), // description
          fc.nat({ max: 100 }), // order
          (title, description, order) => {
            // Mock node response
            const mockNode: GetNodeResponse = {
              node_id: 'test-node-id',
              title: title,
              description: description,
              content: 'Test content',
              markdown_content: '# Test\nSome content',
              boilerplate: '// code here',
              order: order,
              completed: false,
              path_id: 'test-path-id',
            };

            // Mock the hook to return our test node
            vi.mocked(hooks.useNodeDetail).mockReturnValue({
              node: mockNode,
              isLoading: false,
              error: null,
              refetch: vi.fn(),
            });

            const { container } = render(
              <NodeWorkspace pathId="test-path" nodeId="test-node" />
            );

            // Find the title in the header
            const headerTitle = container.querySelector('h1');
            expect(headerTitle).toBeTruthy();
            expect(headerTitle?.textContent).toBe(title);
          }
        ),
        { numRuns: 50 }
      );
    });

    it('should display title consistently across multiple renders with same data', () => {
      fc.assert(
        fc.property(
          fc.string({ minLength: 1, maxLength: 200 }),
          (title) => {
            const mockNode: GetNodeResponse = {
              node_id: 'test-node-id',
              title: title,
              description: 'Test',
              content: 'Test',
              order: 1,
              completed: false,
              path_id: 'test-path-id',
            };

            vi.mocked(hooks.useNodeDetail).mockReturnValue({
              node: mockNode,
              isLoading: false,
              error: null,
              refetch: vi.fn(),
            });

            // Render multiple times
            const { container: container1 } = render(
              <NodeWorkspace pathId="test-path" nodeId="test-node" />
            );
            const title1 = container1.querySelector('h1')?.textContent;

            const { container: container2 } = render(
              <NodeWorkspace pathId="test-path" nodeId="test-node" />
            );
            const title2 = container2.querySelector('h1')?.textContent;

            // Both renders should show the exact same title
            expect(title1).toBe(title);
            expect(title2).toBe(title);
            expect(title1).toBe(title2);
          }
        ),
        { numRuns: 30 }
      );
    });

    it('should handle titles with special HTML characters without escaping issues', () => {
      const titlesWithSpecialChars = [
        'Title with <brackets>',
        'Title with & ampersand',
        'Title with "quotes"',
        "Title with 'apostrophes'",
        'Title with <script>alert("xss")</script>',
        'Title with © symbol',
        'Title with émoji 🚀',
      ];

      titlesWithSpecialChars.forEach(title => {
        const mockNode: GetNodeResponse = {
          node_id: 'test-node-id',
          title: title,
          description: 'Test',
          content: 'Test',
          order: 1,
          completed: false,
          path_id: 'test-path-id',
        };

        vi.mocked(hooks.useNodeDetail).mockReturnValue({
          node: mockNode,
          isLoading: false,
          error: null,
          refetch: vi.fn(),
        });

        const { container } = render(
          <NodeWorkspace pathId="test-path" nodeId="test-node" />
        );

        const headerTitle = container.querySelector('h1');
        expect(headerTitle?.textContent).toBe(title);
      });
    });

    it('should not modify or trim whitespace in titles', () => {
      fc.assert(
        fc.property(
          fc.tuple(
            fc.string({ minLength: 1, maxLength: 50 }),
            fc.nat({ max: 5 }),
            fc.nat({ max: 5 })
          ),
          ([baseTitle, leadingSpaces, trailingSpaces]) => {
            const title = ' '.repeat(leadingSpaces) + baseTitle + ' '.repeat(trailingSpaces);

            const mockNode: GetNodeResponse = {
              node_id: 'test-node-id',
              title: title,
              description: 'Test',
              content: 'Test',
              order: 1,
              completed: false,
              path_id: 'test-path-id',
            };

            vi.mocked(hooks.useNodeDetail).mockReturnValue({
              node: mockNode,
              isLoading: false,
              error: null,
              refetch: vi.fn(),
            });

            const { container } = render(
              <NodeWorkspace pathId="test-path" nodeId="test-node" />
            );

            const headerTitle = container.querySelector('h1');
            expect(headerTitle?.textContent).toBe(title);
          }
        ),
        { numRuns: 30 }
      );
    });

    it('should display correct order number from API response', () => {
      fc.assert(
        fc.property(
          fc.nat({ max: 1000 }),
          fc.string({ minLength: 1, maxLength: 100 }),
          (order, title) => {
            const mockNode: GetNodeResponse = {
              node_id: 'test-node-id',
              title: title,
              description: 'Test',
              content: 'Test',
              order: order,
              completed: false,
              path_id: 'test-path-id',
            };

            vi.mocked(hooks.useNodeDetail).mockReturnValue({
              node: mockNode,
              isLoading: false,
              error: null,
              refetch: vi.fn(),
            });

            const { container } = render(
              <NodeWorkspace pathId="test-path" nodeId="test-node" />
            );

            const orderText = container.textContent || '';
            expect(orderText).toContain(`Ejercicio ${order}`);
          }
        ),
        { numRuns: 50 }
      );
    });
  });
});

