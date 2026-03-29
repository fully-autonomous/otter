// Noms TUI - Full Interactive TUI

import { getLog, getBranches, getStatus, getDatasets, checkoutBranch, createBranch, deleteBranch, getRemotes, push, pull, addRemote, commit, runNoms } from "./lib/noms";
import type { NomsLogEntry, NomsBranch, NomsStatus, NomsDataset, Remote } from "./lib/noms";

type Screen = 'dashboard' | 'branches' | 'datasets' | 'dataset-view' | 'sync' | 'commit' | 'commit-message' | 'branch-create' | 'branch-delete' | 'add-remote';

const state = {
  screen: 'dashboard' as Screen,
  branches: [] as NomsBranch[],
  status: null as NomsStatus | null,
  datasets: [] as NomsDataset[],
  commits: [] as NomsLogEntry[],
  remotes: [] as Remote[],
  selectedIndex: 0,
  inputBuffer: '',
  inputBuffer2: '', // For second field (e.g., URL in add remote)
  datasetContent: '', // For viewing dataset details
  commitPath: '',
  commitDataset: '',
};

function header(title: string): string {
  return `\x1b[1;36m╔════════════════════════════════════════════════════════╗\n` +
         `║              Noms Terminal UI - ${title}${' '.repeat(39 - title.length)}║\n` +
         `╚════════════════════════════════════════════════════════╝\x1b[0m\n\n`;
}

function renderDashboard(): string {
  let output = header('Dashboard');
  
  if (!state.status && state.branches.length === 0) {
    output += `\x1b[33m⚠ Not initialized. Run "noms init" first.\x1b[0m\n\n`;
  }
  
  output += `\x1b[1m┌─ Status\x1b[0m─────────────────────────────────────────────┐\n`;
  if (state.status) {
    output += `│ \x1b[32m✓ Working directory clean\x1b[0m\n`;
    output += `│ Branch: \x1b[36m${state.status.current_branch}\x1b[0m\n`;
  } else if (state.branches.length > 0) {
    output += `│ \x1b[33m○ No working directory\x1b[0m\n`;
  }
  output += `│ Datasets: ${state.datasets.length}\n`;
  output += `│ Remotes: ${state.remotes.length}\n`;
  output += `└────────────────────────────────────────────────────────┘\n\n`;
  
  output += `\x1b[1m┌─ Branches\x1b[0m──────────────────────────────────────────────┐\n`;
  if (state.branches.length === 0) {
    output += `│ (no branches)\x1b[2m                                            │\x1b[0m\n`;
  } else {
    state.branches.slice(0, 5).forEach(branch => {
      const isCurrent = state.status?.current_branch === branch.name;
      const marker = isCurrent ? `\x1b[32m●\x1b[0m` : ` `;
      const isSelected = state.branches.indexOf(branch) === state.selectedIndex;
      const sel = isSelected ? `\x1b[44m \x1b[0m` : `  `;
      output += `│${sel}${marker} ${branch.name}\n`;
    });
    if (state.branches.length > 5) {
      output += `│   ... and ${state.branches.length - 5} more\x1b[2m                             │\x1b[0m\n`;
    }
  }
  output += `└────────────────────────────────────────────────────────┘\n\n`;
  
  output += `\x1b[1m┌─ Recent Commits\x1b[0m────────────────────────────────────────┐\n`;
  if (state.commits.length === 0) {
    output += `│ (no commits yet)\x1b[2m                                         │\x1b[0m\n`;
  } else {
    state.commits.slice(0, 5).forEach(commit => {
      const shortHash = commit.hash.substring(0, 7);
      const msg = commit.message.length > 35 ? commit.message.substring(0, 32) + '...' : commit.message;
      output += `│ \x1b[2m${shortHash}\x1b[0m ${msg}\n`;
    });
  }
  output += `└────────────────────────────────────────────────────────┘\n\n`;
  
  output += `\x1b[2m  [\x1b[0m1\x1b[2m] Dashboard    [\x1b[0m2\x1b[2m] Branches    [\x1b[0m3\x1b[2m] Datasets    [\x1b[0m4\x1b[2m] Sync    [\x1b[0m5\x1b[2m] Commit    [\x1b[0mq\x1b[2m] Quit\x1b[0m\n`;
  
  return output;
}

function renderBranches(): string {
  let output = header('Branches');
  
  output += `\x1b[1m┌─ Your Branches\x1b[0m────────────────────────────────────────────┐\n`;
  if (state.branches.length === 0) {
    output += `│ (no branches yet)\x1b[2m                                          │\x1b[0m\n`;
  } else {
    state.branches.forEach((branch, i) => {
      const isCurrent = state.status?.current_branch === branch.name;
      const marker = isCurrent ? `\x1b[32m●\x1b[0m` : ` `;
      const isSelected = i === state.selectedIndex;
      const sel = isSelected ? `\x1b[44m \x1b[0m` : `  `;
      const name = isCurrent ? `\x1b[1m${branch.name}\x1b[0m` : branch.name;
      output += `│${sel}${marker} ${name}\n`;
    });
  }
  output += `└────────────────────────────────────────────────────────┘\n\n`;
  
  output += `\x1b[1mCommands:\x1b[0m\n`;
  output += `  [\x1b[0mn\x1b[2m] new branch    [\x1b[0mc\x1b[2m] checkout    [\x1b[0md\x1b[2m] delete    [\x1b[0mb\x1b[2m] back to dashboard\x1b[0m\n`;
  
  return output;
}

function renderBranchCreate(): string {
  let output = header('New Branch');
  
  output += `  Enter new branch name:\n\n`;
  output += `  \x1b[36m> ${state.inputBuffer}\x1b[0m\x1b[36m▌\x1b[0m\n\n`;
  output += `  \x1b[2m[Enter] create    [Escape] cancel\x1b[0m\n`;
  
  return output;
}

function renderBranchDelete(): string {
  const branch = state.branches[state.selectedIndex];
  let output = header('Delete Branch');
  
  output += `  \x1b[31m⚠ Delete branch:\x1b[0m \x1b[33m${branch?.name || 'none'}\x1b[0m\n\n`;
  output += `  \x1b[2m[Enter] confirm delete    [Escape] cancel\x1b[0m\n`;
  
  return output;
}

function renderDatasets(): string {
  let output = header('Datasets');
  
  output += `\x1b[1m┌─ Datasets in this database\x1b[0m─────────────────────────────────────┐\n`;
  if (state.datasets.length === 0) {
    output += `│ (no datasets yet)\x1b[2m                                             │\x1b[0m\n`;
  } else {
    state.datasets.forEach((ds, i) => {
      const isSelected = i === state.selectedIndex;
      const sel = isSelected ? `\x1b[44m \x1b[0m` : `  `;
      const head = ds.head.substring(0, 12);
      output += `│${sel}${ds.name.padEnd(30)} ${head}\n`;
    });
  }
  output += `└────────────────────────────────────────────────────────┘\n\n`;
  
  output += `  [\x1b[0mEnter\x1b[2m] view dataset    [\x1b[0mb\x1b[2m] back to dashboard\x1b[0m\n`;
  
  return output;
}

function renderDatasetView(): string {
  const ds = state.datasets[state.selectedIndex];
  let output = header('Dataset View');
  
  if (!ds) {
    output += `  \x1b[33mNo dataset selected\x1b[0m\n\n`;
    output += `  [\x1b[0mb\x1b[2m] back\x1b[0m\n`;
    return output;
  }
  
  output += `  \x1b[1mName:\x1b[0m ${ds.name}\n`;
  output += `  \x1b[1mHead:\x1b[0m ${ds.head}\n\n`;
  
  output += `\x1b[1m┌─ Content\x1b[0m────────────────────────────────────────────────────┐\n`;
  const lines = state.datasetContent.split('\n').slice(0, 20);
  if (lines.length === 0 || state.datasetContent === '') {
    output += `│ (loading or empty)\x1b[2m                                           │\x1b[0m\n`;
  } else {
    lines.forEach(line => {
      const truncated = line.length > 74 ? line.substring(0, 71) + '...' : line;
      output += `│ ${truncated.padEnd(74)}\n`;
    });
    if (state.datasetContent.split('\n').length > 20) {
      output += `│ ... (more content)\x1b[2m                                           │\x1b[0m\n`;
    }
  }
  output += `└────────────────────────────────────────────────────────┘\n\n`;
  
  output += `  [\x1b[0mb\x1b[2m] back to datasets    [\x1b[0mr\x1b[2m] refresh\x1b[0m\n`;
  
  return output;
}

function renderSync(): string {
  let output = header('Sync');
  
  output += `\x1b[1m┌─ Remotes\x1b[0m───────────────────────────────────────────────────┐\n`;
  if (state.remotes.length === 0) {
    output += `│ (no remotes configured)\x1b[2m                                       │\x1b[0m\n`;
    output += `│ Add a remote to sync your database\x1b[2m                               │\x1b[0m\n`;
  } else {
    state.remotes.forEach((remote, i) => {
      const isSelected = i === state.selectedIndex;
      const sel = isSelected ? `\x1b[44m \x1b[0m` : `  `;
      output += `│${sel}\x1b[36m${remote.name}\x1b[0m → ${remote.url}\n`;
    });
  }
  output += `└────────────────────────────────────────────────────────┘\n\n`;
  
  output += `\x1b[1mActions:\x1b[0m\n`;
  output += `  [\x1b[0mp\x1b[2m] push to remote    [\x1b[0ml\x1b[2m] pull from remote\n`;
  output += `  [\x1b[0ma\x1b[2m] add remote\n`;
  output += `  [\x1b[0mb\x1b[2m] back to dashboard\x1b[0m\n`;
  
  return output;
}

function renderCommit(): string {
  let output = header('Commit');
  
  if (!state.status) {
    output += `  \x1b[33m⚠ Not initialized\x1b[0m\n`;
    output += `  Run "noms init" first\n\n`;
    output += `  [\x1b[0mb\x1b[2m] back to dashboard\x1b[0m\n`;
    return output;
  }
  
  output += `  Branch: \x1b[36m${state.status.current_branch}\x1b[0m\n\n`;
  
  if (state.status.modified.length === 0 && state.status.staged.length === 0) {
    output += `  \x1b[32m✓ Working directory is clean\x1b[0m\n`;
    output += `  No changes to commit\n`;
  } else {
    if (state.status.modified.length > 0) {
      output += `  \x1b[33mModified files:\x1b[0m\n`;
      state.status.modified.forEach(f => {
        output += `    \x1b[33mM\x1b[0m ${f}\n`;
      });
      output += `\n`;
    }
    
    if (state.status.staged.length > 0) {
      output += `  \x1b[32mStaged files:\x1b[0m\n`;
      state.status.staged.forEach(f => {
        output += `    \x1b[32mA\x1b[0m ${f}\n`;
      });
      output += `\n`;
    }
  }
  
  output += `\x1b[1m─\x1b[0m`.repeat(56) + `\n\n`;
  
  output += `  [\x1b[0mc\x1b[2m] write commit message\n`;
  output += `  [\x1b[0mb\x1b[2m] back to dashboard\x1b[0m\n`;
  
  return output;
}

function renderCommitMessage(): string {
  let output = header('Write Commit');
  
  output += `  Enter commit details:\n\n`;
  output += `  Message: \x1b[36m${state.inputBuffer}\x1b[0m\x1b[36m▌\x1b[0m\n`;
  output += `  Path:    \x1b[36m${state.commitPath}\x1b[0m\n`;
  output += `  Dataset: \x1b[36m${state.commitDataset}\x1b[0m\n\n`;
  
  output += `  \x1b[2m[Tab] next field  [Enter] commit  [Escape] cancel\x1b[0m\n`;
  
  return output;
}

function renderAddRemote(): string {
  let output = header('Add Remote');
  
  output += `  Enter remote details:\n\n`;
  output += `  Name: \x1b[36m${state.inputBuffer}\x1b[0m\x1b[36m▌\x1b[0m\n`;
  output += `  URL:  \x1b[36m${state.inputBuffer2}\x1b[0m\n\n`;
  
  output += `  \x1b[2m[Tab] next field  [Enter] add  [Escape] cancel\x1b[0m\n`;
  
  return output;
}

async function loadData() {
  try {
    const [status, branches, datasets, commits, remotes] = await Promise.all([
      getStatus(),
      getBranches(),
      getDatasets(),
      getLog(10),
      getRemotes(),
    ]);
    state.status = status;
    state.branches = branches;
    state.datasets = datasets;
    state.commits = commits;
    state.remotes = remotes;
  } catch (e) {
    // Ignore errors - db might not be initialized
  }
}

function render(): string {
  switch (state.screen) {
    case 'dashboard': return renderDashboard();
    case 'branches': return renderBranches();
    case 'branch-create': return renderBranchCreate();
    case 'branch-delete': return renderBranchDelete();
    case 'datasets': return renderDatasets();
    case 'dataset-view': return renderDatasetView();
    case 'sync': return renderSync();
    case 'add-remote': return renderAddRemote();
    case 'commit': return renderCommit();
    case 'commit-message': return renderCommitMessage();
    default: return renderDashboard();
  }
}

async function handleInput(input: string): Promise<void> {
  const key = input.trim().toLowerCase();
  
  // Handle branch create mode
  if (state.screen === 'branch-create') {
    if (key === '') {
      // Enter pressed - create branch
      if (state.inputBuffer.trim()) {
        await createBranch(state.inputBuffer.trim());
        await loadData();
      }
      state.screen = 'branches';
      state.inputBuffer = '';
    } else if (input === '\x1b') {
      // Escape
      state.screen = 'branches';
      state.inputBuffer = '';
    } else if (input === '\x7f') {
      // Backspace
      state.inputBuffer = state.inputBuffer.slice(0, -1);
    } else if (input.length === 1) {
      state.inputBuffer += input;
    }
    return;
  }
  
  // Handle branch delete mode
  if (state.screen === 'branch-delete') {
    if (key === '') {
      // Enter pressed - delete branch
      const branch = state.branches[state.selectedIndex];
      if (branch && branch.name !== state.status?.current_branch) {
        await deleteBranch(branch.name);
        await loadData();
      }
      state.screen = 'branches';
    } else if (input === '\x1b') {
      state.screen = 'branches';
    }
    return;
  }
  
  // Handle add remote mode
  if (state.screen === 'add-remote') {
    if (key === '') {
      // Enter pressed - add remote
      const name = state.inputBuffer.trim();
      const url = state.inputBuffer2.trim();
      if (name && url) {
        await addRemote(name, url);
        await loadData();
        console.log(`\n\x1b[32mAdded remote: ${name}\x1b[0m\n`);
      }
      state.screen = 'sync';
      state.inputBuffer = '';
      state.inputBuffer2 = '';
    } else if (input === '\x1b') {
      state.screen = 'sync';
      state.inputBuffer = '';
      state.inputBuffer2 = '';
    } else if (input === '\t') {
      // Tab - switch between fields
      // For simplicity, just track which field is active
    } else if (input === '\x7f') {
      // Backspace
      state.inputBuffer = state.inputBuffer.slice(0, -1);
    } else if (input.length === 1) {
      state.inputBuffer += input;
    }
    return;
  }
  
  // Handle commit message mode
  if (state.screen === 'commit-message') {
    if (key === '') {
      // Enter pressed - commit
      const message = state.inputBuffer.trim();
      const path = state.commitPath.trim() || '.';
      const dataset = state.commitDataset.trim() || 'main';
      if (message) {
        await commit(message, path, dataset);
        await loadData();
        console.log(`\n\x1b[32mCommitted: ${message}\x1b[0m\n`);
      }
      state.screen = 'commit';
      state.inputBuffer = '';
      state.commitPath = '';
      state.commitDataset = '';
    } else if (input === '\x1b') {
      state.screen = 'commit';
      state.inputBuffer = '';
      state.commitPath = '';
      state.commitDataset = '';
    } else if (input === '\t') {
      // Tab - could cycle through fields in future
    } else if (input === '\x7f') {
      // Backspace
      state.inputBuffer = state.inputBuffer.slice(0, -1);
    } else if (input.length === 1) {
      state.inputBuffer += input;
    }
    return;
  }
  
  // Screen navigation
  switch (state.screen) {
    case 'dashboard':
      if (key === '1' || key === 'd') state.screen = 'dashboard';
      else if (key === '2' || key === 'b') { state.screen = 'branches'; state.selectedIndex = 0; }
      else if (key === '3') { state.screen = 'datasets'; state.selectedIndex = 0; }
      else if (key === '4' || key === 's') { state.screen = 'sync'; state.selectedIndex = 0; }
      else if (key === '5' || key === 'c') state.screen = 'commit';
      else if (key === 'r') await loadData();
      else if (key === 'q') { console.log('\n\nThanks for using Noms! 👋\n'); process.exit(0); }
      break;
      
    case 'branches':
      if (key === 'n') state.screen = 'branch-create';
      else if (key === 'c' && state.branches[state.selectedIndex]) {
        await checkoutBranch(state.branches[state.selectedIndex].name);
        await loadData();
      } else if (key === 'd' && state.branches[state.selectedIndex]) {
        state.screen = 'branch-delete';
      } else if (key === 'b' || key === 'q' || key === '1') state.screen = 'dashboard';
      else if (key === 'arrowup' || key === 'k') state.selectedIndex = Math.max(0, state.selectedIndex - 1);
      else if (key === 'arrowdown' || key === 'j') state.selectedIndex = Math.min(state.branches.length - 1, state.selectedIndex + 1);
      break;
      
    case 'datasets':
      if (key === 'b' || key === 'q' || key === '1') state.screen = 'dashboard';
      else if (key === 'arrowup' || key === 'k') state.selectedIndex = Math.max(0, state.selectedIndex - 1);
      else if (key === 'arrowdown' || key === 'j') state.selectedIndex = Math.min(state.datasets.length - 1, state.selectedIndex + 1);
      else if (key === '') {
        // Enter - view dataset
        const ds = state.datasets[state.selectedIndex];
        if (ds) {
          state.screen = 'dataset-view';
          // Load dataset content
          const result = await runNoms(['show', ds.name], 'text');
          state.datasetContent = result.stdout || '(empty)';
        }
      }
      break;
      
    case 'dataset-view':
      if (key === 'b' || key === 'q' || key === '1' || key === '') {
        state.screen = 'datasets';
        state.datasetContent = '';
      } else if (key === 'r') {
        // Refresh
        const ds = state.datasets[state.selectedIndex];
        if (ds) {
          const result = await runNoms(['show', ds.name], 'text');
          state.datasetContent = result.stdout || '(empty)';
        }
      }
      break;
      
    case 'sync':
      if (key === 'a') {
        state.screen = 'add-remote';
        state.inputBuffer = '';
        state.inputBuffer2 = '';
      } else if (key === 'p' && state.remotes[state.selectedIndex]) {
        await push(state.remotes[state.selectedIndex].name);
        console.log('\n\x1b[32mPush complete!\x1b[0m\n');
      } else if (key === 'l' && state.remotes[state.selectedIndex]) {
        await pull(state.remotes[state.selectedIndex].name);
        console.log('\n\x1b[32mPull complete!\x1b[0m\n');
        await loadData();
      } else if (key === 'b' || key === 'q' || key === '1') state.screen = 'dashboard';
      else if (key === 'arrowup' || key === 'k') state.selectedIndex = Math.max(0, state.selectedIndex - 1);
      else if (key === 'arrowdown' || key === 'j') state.selectedIndex = Math.min(state.remotes.length - 1, state.selectedIndex + 1);
      break;
      
    case 'commit':
      if (key === 'c') {
        state.screen = 'commit-message';
        state.inputBuffer = '';
        state.commitPath = '';
        state.commitDataset = '';
      } else if (key === 'b' || key === 'q' || key === '1') state.screen = 'dashboard';
      break;
  }
}

async function main() {
  await loadData();
  
  console.clear();
  console.log(render());
  
  const readline = await import('readline');
  
  function ask() {
    const rl = readline.createInterface({
      input: process.stdin,
      output: process.stdout,
    });
    
    rl.question('', async (input) => {
      rl.close();
      await handleInput(input);
      console.clear();
      console.log(render());
      ask();
    });
  }
  
  ask();
}

main().catch(console.error);
