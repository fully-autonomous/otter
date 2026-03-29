// Noms CLI wrapper library
// Communicates with the noms CLI via subprocess

export interface NomsLogEntry {
  hash: string;
  parents: string;
  message: string;
  author: string;
  date: string;
  dataset: string;
}

export interface NomsBranch {
  name: string;
  head: string;
  created_at: string;
}

export interface NomsDataset {
  name: string;
  head: string;
  root_hash: string;
}

export interface NomsStatus {
  current_branch: string;
  branches: string[];
  staged: string[];
  modified: string[];
  clean: boolean;
}

export interface Remote {
  name: string;
  url: string;
}

export interface NomsError {
  error: string;
}

function isNomsError(obj: any): obj is NomsError {
  return obj && obj.error;
}

export async function runNoms(args: string[], format: 'text' | 'json' = 'json'): Promise<{ stdout: string; stderr: string; code: number }> {
  const proc = Bun.spawn(['noms', ...args, '--format', format], {
    stdout: 'pipe',
    stderr: 'pipe',
  });
  
  const [stdout, stderr] = await Promise.all([
    new Response(proc.stdout).text(),
    new Response(proc.stderr).text(),
  ]);
  
  const code = await proc.exited;
  return { stdout, stderr, code };
}

export async function getLog(limit: number = 50): Promise<NomsLogEntry[]> {
  const result = await runNoms(['query', 'noms_log']);
  if (result.code !== 0) {
    throw new Error(result.stderr || result.stdout);
  }
  
  try {
    const data = JSON.parse(result.stdout);
    return data.noms_log || [];
  } catch {
    return [];
  }
}

export async function getBranches(): Promise<NomsBranch[]> {
  const result = await runNoms(['query', 'noms_branches']);
  if (result.code !== 0) {
    throw new Error(result.stderr);
  }
  
  try {
    const data = JSON.parse(result.stdout);
    return data.branches || [];
  } catch {
    return [];
  }
}

export async function getDatasets(): Promise<NomsDataset[]> {
  const result = await runNoms(['query', 'noms_datasets']);
  if (result.code !== 0) {
    throw new Error(result.stderr);
  }
  
  try {
    const data = JSON.parse(result.stdout);
    return data.datasets || [];
  } catch {
    return [];
  }
}

export async function getStatus(): Promise<NomsStatus | null> {
  const result = await runNoms(['status']);
  if (result.code !== 0) {
    return null;
  }
  
  try {
    return JSON.parse(result.stdout);
  } catch {
    return null;
  }
}

export async function getRemotes(): Promise<Remote[]> {
  const result = await runNoms(['remote']);
  if (result.code !== 0) {
    return [];
  }
  
  try {
    return JSON.parse(result.stdout);
  } catch {
    return [];
  }
}

export async function checkoutBranch(branchName: string, create: boolean = false): Promise<boolean> {
  const args = create ? ['checkout', '-b', branchName] : ['checkout', branchName];
  const result = await runNoms(args);
  return result.code === 0;
}

export async function createBranch(branchName: string): Promise<boolean> {
  const result = await runNoms(['branch', '-c', branchName]);
  return result.code === 0;
}

export async function deleteBranch(branchName: string): Promise<boolean> {
  const result = await runNoms(['branch', '-d', branchName]);
  return result.code === 0;
}

export async function addRemote(name: string, url: string): Promise<boolean> {
  const result = await runNoms(['remote', '--add', name, url]);
  return result.code === 0;
}

export async function removeRemote(name: string): Promise<boolean> {
  const result = await runNoms(['remote', '--remove', name]);
  return result.code === 0;
}

export async function push(remote: string = 'origin', branch: string = ''): Promise<boolean> {
  const args = branch ? ['push', remote, branch] : ['push', remote];
  const result = await runNoms(args);
  return result.code === 0;
}

export async function pull(remote: string = 'origin', branch: string = ''): Promise<boolean> {
  const args = branch ? ['pull', remote, branch] : ['pull', remote];
  const result = await runNoms(args);
  return result.code === 0;
}

export async function init(): Promise<boolean> {
  const result = await runNoms(['init']);
  return result.code === 0;
}

export async function commit(message: string, path: string, dataset: string): Promise<boolean> {
  const result = await runNoms(['commit', '-m', message, path, dataset]);
  return result.code === 0;
}
