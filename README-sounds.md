# WC3 Sounds for Claude Code

This configuration adds Warcraft III character sounds to various Claude Code
events, making your coding experience more immersive!

## Sound Source

All sounds are from the HiveWorkshop community:
https://www.hiveworkshop.com/threads/sound-sets-from-campaign-dialog.328584/

## Quick Setup

Run the setup script to automatically download and install the sounds:

```bash
~/.claude/scripts/setup-sounds.sh
```

## Manual Setup

If the automatic setup fails, you can manually install the sounds:

1. Download the sound sets from:
   https://www.hiveworkshop.com/threads/sound-sets-from-campaign-dialog.328584/
2. Extract the ZIP file to `~/Sounds/game_samples/`
3. Ensure the directory structure looks like:
   ```
   ~/Sounds/game_samples/
   ├── Human/
   │   ├── Arthas/
   │   ├── Footman/
   │   ├── Knight/
   │   └── ...
   ├── Orc/
   │   ├── Grunt/
   │   ├── Peon/
   │   └── ...
   ├── NightElf/
   │   ├── Illidan/
   │   ├── Furion/
   │   └── ...
   └── jobs-done_1.mp3
   ```

## Sound Mappings

### Claude Events (settings.json)

- **Edit/Write files**: Arthas "Yes" - acknowledgment of orders
- **Bash commands**: "Jobs Done!" - task completion
- **Read files**: Illidan "What?" - seeking knowledge
- **Search (Grep/Glob/Task)**: Rexxar "What?" - hunting for information
- **TodoWrite**: Villager "What?" - work assignment
- **Web operations**: Medivh "What?" - mystical knowledge seeking
- **ExitPlanMode**: Footman Warcry - ready for action
- **UserPromptSubmit**: Paladin "Ready!" - beginning a new quest
- **Stop**: Furion "Ready" - task complete
- **Error**: Arthas "Pissed" - frustration with errors

### File Type Specific (post-edit.sh)

- **Go files**: Orc Grunt "Yes" - strong and reliable
- **Markdown**: Night Elf Tyrande "Yes" - elegant documentation
- **JS/TS/JSON**: Human Knight "Yes" - noble web development
- **YAML**: Medivh "Yes" - configuration wisdom
- **Rust**: Rexxar "Yes" - fierce and powerful

### Script Events

- **validate-claude.sh success**: Cairne "Yes" - wise approval
- **validate-claude.sh failure**: Garithos "Pissed" - validation anger
- **qsweep.sh start**: Furion "What?" - beginning the sweep
- **qsweep.sh no files**: Peon "Yes" - work complete
- **qsweep.sh success**: Paladin Warcry - triumphant completion

## Troubleshooting

### Sounds not playing?

1. Check if `afplay` is available: `which afplay`
2. Verify sounds are installed: `ls ~/Sounds/game_samples/`
3. Test a sound manually: `afplay ~/Sounds/game_samples/jobs-done_1.mp3`

### Download fails?

- The HiveWorkshop might require login or have changed URLs
- Download manually from the link above and extract to `~/Sounds/game_samples/`

## Customization

To change sounds, edit `~/.claude/settings.json` and update the sound paths in
the hooks section. Each hook can play any sound file using:

```json
{
  "type": "command",
  "command": "nohup afplay ~/Sounds/game_samples/YOUR_SOUND.mp3 > /dev/null 2>&1 &"
}
```

## Credits

All sounds are property of Blizzard Entertainment and extracted by the WC3
community at HiveWorkshop.
