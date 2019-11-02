import * as monaco from 'monaco-editor'
import mylang from './mylang.js'

self.MonacoEnvironment = {
  getWorkerUrl: function(moduleId, label) {
    if (label === 'json') {
      return './json.worker.js';
    }
    if (label === 'css') {
      return './css.worker.js';
    }
    if (label === 'html') {
      return './html.worker.js';
    }
    if (label === 'typescript' || label === 'javascript') {
      return './ts.worker.js';
    }
    return './editor.worker.js';
  },
};

monaco.languages.register({
  id: 'mylang'
});
monaco.languages.setMonarchTokensProvider('mylang', {
  keywords: [
    "func", "class", "for", "if", "return", "break",
  ],
  typeKeywords: [
    "int", "error",
  ],
  tokenizer: {
    root: [
      [/[a-z_]\w*/, {
        cases: {
          "@keywords": "keyword",
          "@typeKeywords": "type",
          "@default": "variable",
        },
      }],
      [/\(\*/, 'comment', '@comment'],
      [/"/, "string", "@string"],
    ],
    string: [
      [/[^"]+/, "string"],
      [/"/, "string", "@pop"],
    ],
    comment: [
      [/[^*(]+/, 'comment'],
      [/\(\*+/, 'comment', '@comment'],
      [/\*[^)]/, 'comment'],
      [/\*\)/, 'comment', '@pop'],
    ],
  }
});

monaco.languages.registerCompletionItemProvider('mylang', {
  provideCompletionItems: function(model, pos, context, token) {
    let suggestions = []
    suggestions.push({
      insertText: "String",
      label: "String",
      kind: monaco.languages.CompletionItemKind.Function,
    });
    suggestions.push({
      insertText: "Abs",
      label: "Abs",
      kind: monaco.languages.CompletionItemKind.Function,
    });
    return {
      suggestions: suggestions,
    }
  }
})

monaco.editor.create(document.getElementById('editor'), {
  language: 'mylang',
  theme: 'vs',
});

