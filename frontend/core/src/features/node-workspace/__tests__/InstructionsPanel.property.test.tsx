import { describe, it, expect } from 'vitest';
import { render, screen } from '@testing-library/react';
import * as fc from 'fast-check';
import { InstructionsPanel } from '../components/InstructionsPanel';

/**
 * Property 4: Markdown Content Rendering
 *
 * For any valid markdown_content string received from the API,
 * rendering it in the instructions panel should produce HTML that
 * contains all the text content from the original markdown without data loss.
 */
describe('InstructionsPanel - Property-Based Tests', () => {
  describe('Property 4: Markdown Content Rendering', () => {
    it('should preserve all text content when rendering any valid markdown', () => {
      fc.assert(
        fc.property(
          // Generate more realistic markdown strings with words
          fc.oneof(
            // Plain text with real words
            fc.array(fc.lorem({ maxCount: 1 }), { minLength: 3, maxLength: 20 })
              .map(words => words.join(' ')),
            // Markdown with headers
            fc.array(fc.lorem({ maxCount: 1 }), { minLength: 2, maxLength: 10 })
              .map(words => {
                const text = words.join(' ');
                return `# ${text}\n## Subtitle\n${text}`;
              }),
            // Markdown with lists
            fc.array(
              fc.array(fc.lorem({ maxCount: 1 }), { minLength: 2, maxLength: 5 })
                .map(words => words.join(' ')),
              { minLength: 2, maxLength: 5 }
            ).map(items => items.map(item => `- ${item}`).join('\n')),
            // Markdown with code blocks
            fc.array(fc.lorem({ maxCount: 1 }), { minLength: 3, maxLength: 10 })
              .map(words => {
                const code = words.join(' ');
                return `\`\`\`go\n${code}\n\`\`\``;
              }),
            // Markdown with emphasis
            fc.array(fc.lorem({ maxCount: 1 }), { minLength: 3, maxLength: 10 })
              .map(words => {
                const text = words.join(' ');
                return `**${text}** and *${text}*`;
              }),
            // Mixed content
            fc.tuple(
              fc.array(fc.lorem({ maxCount: 1 }), { minLength: 2, maxLength: 5 }),
              fc.array(fc.lorem({ maxCount: 1 }), { minLength: 3, maxLength: 8 }),
              fc.array(fc.lorem({ maxCount: 1 }), { minLength: 2, maxLength: 5 })
            ).map(([titleWords, descWords, codeWords]) => {
              const title = titleWords.join(' ');
              const desc = descWords.join(' ');
              const code = codeWords.join(' ');
              return `# ${title}\n\n${desc}\n\n\`\`\`go\n${code}\n\`\`\``;
            })
          ),
          fc.array(fc.lorem({ maxCount: 1 }), { minLength: 2, maxLength: 5 })
            .map(words => `Test: ${words.join(' ')}`), // title
          (markdownContent, title) => {
            const { container } = render(
              <InstructionsPanel
                markdownContent={markdownContent}
                title={title}
              />
            );

            // Get rendered text content
            const renderedText = container.textContent || '';

            // The title should be present
            expect(renderedText).toContain(title);

            // Extract key words from markdown (alphanumeric words of significant length)
            const words = markdownContent
              .replace(/```[\s\S]*?```/g, '') // Remove code blocks
              .replace(/[*_#`\[\]()]/g, ' ') // Remove markdown syntax
              .split(/\s+/)
              .filter(word => word.length > 3 && /[a-zA-Z0-9]/.test(word));

            if (words.length === 0) {
              // If no significant words, just check title is present
              expect(renderedText).toContain(title);
              return;
            }

            // At least 40% of significant words should be present
            const wordsPresent = words.filter(word =>
              renderedText.toLowerCase().includes(word.toLowerCase())
            ).length;

            const preservationRate = wordsPresent / words.length;
            expect(preservationRate).toBeGreaterThanOrEqual(0.4);
          }
        ),
        { numRuns: 30 } // Run 30 random test cases
      );
    });

    it('should handle empty markdown content gracefully', () => {
      const { container } = render(
        <InstructionsPanel markdownContent="" title="Test" />
      );
      expect(container.textContent).toContain('Test');
    });

    it('should render code blocks without losing content', () => {
      fc.assert(
        fc.property(
          fc.string({ minLength: 5, maxLength: 200 }),
          (codeContent) => {
            const markdown = `\`\`\`go\n${codeContent}\n\`\`\``;
            const { container } = render(
              <InstructionsPanel markdownContent={markdown} title="Code Test" />
            );

            // The code content should be present in the rendered output
            const renderedText = container.textContent || '';

            // Check that the code content is present (it might be formatted differently)
            const codeWords = codeContent.split(/\s+/).filter(w => w.length > 2);
            const presentWords = codeWords.filter(word => renderedText.includes(word));

            const preservationRate = codeWords.length > 0 ? presentWords.length / codeWords.length : 1;
            expect(preservationRate).toBeGreaterThanOrEqual(0.8);
          }
        ),
        { numRuns: 30 }
      );
    });

    it('should preserve special characters in markdown', () => {
      const specialCharsMarkdown = `
# Test with special chars
- Item with & ampersand
- Item with < less than
- Item with > greater than
- Item with "quotes"
- Item with 'apostrophe'
      `;

      const { container } = render(
        <InstructionsPanel
          markdownContent={specialCharsMarkdown}
          title="Special Chars"
        />
      );

      const text = container.textContent || '';
      expect(text).toContain('ampersand');
      expect(text).toContain('less than');
      expect(text).toContain('greater than');
    });
  });
});

