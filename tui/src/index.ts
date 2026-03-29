// Noms TUI - Main Entry Point
// This is a template that shows the structure for an OpenTUI-based TUI

import { createCliRenderer, Box, Text, Flex, Input } from "@opentui/core";

async function main() {
  const renderer = await createCliRenderer();
  
  console.clear();
  
  // Welcome screen
  const box = new Box(renderer, {
    width: "100%",
    height: "100%",
    borderStyle: "round",
    title: "Noms TUI",
  });
  
  box.render();
  
  // Show welcome message
  console.log(`
╔═══════════════════════════════════════════════════════════╗
║                    Noms Terminal UI                       ║
╠═══════════════════════════════════════════════════════════╣
║                                                           ║
║  Welcome to the Noms TUI!                                ║
║                                                           ║
║  This is a template for the OpenTUI-based interface.     ║
║                                                           ║
║  Current CLI commands available:                          ║
║    - noms init          Initialize database                ║
║    - noms branch       Manage branches                    ║
║    - noms checkout     Switch branches                     ║
║    - noms status       Show status                        ║
║    - noms log         View commit history                 ║
║    - noms remote      Manage remotes                      ║
║    - noms push        Push to remote                      ║
║    - noms pull        Pull from remote                    ║
║    - noms clone       Clone remote                        ║
║                                                           ║
║  To launch the full TUI, run:                            ║
║    bun run src/index.ts                                   ║
║                                                           ║
║  Note: This requires OpenTUI to be fully integrated.     ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
  `);
  
  // For now, just display a message
  // Full OpenTUI integration would require:
  // 1. Installing @opentui/core with bun
  // 2. Building interactive components
  // 3. Setting up event handlers
  
  process.exit(0);
}

main().catch(console.error);
