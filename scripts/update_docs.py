#!/usr/bin/env python3
"""
Daily documentation refresh script using Anthropic API directly.
More reliable alternative to using Claude CLI in headless mode.
"""

import os
import sys
import json
import requests
from datetime import datetime
from pathlib import Path
from typing import List, Dict
import re

# Try to import anthropic, provide helpful message if not installed
try:
    import anthropic
except ImportError:
    print("❌ anthropic package not installed.")
    print("   Run: pip install anthropic")
    sys.exit(1)

def slugify(url: str) -> str:
    """Convert URL to a safe filename."""
    # Remove protocol and domain
    path = re.sub(r'https?://[^/]+/', '', url)
    # Replace slashes with dashes
    path = path.replace('/', '-')
    # Remove any remaining special characters
    path = re.sub(r'[^\w\-]', '', path)
    # Add .md extension if not present
    if not path.endswith('.md'):
        path += '.md'
    return path

def fetch_with_anthropic(urls: List[str], docs_dir: str = "docs") -> Dict:
    """Fetch documentation using Anthropic API."""
    
    # Get API key from environment
    api_key = os.environ.get('ANTHROPIC_API_KEY')
    if not api_key:
        print("❌ ANTHROPIC_API_KEY environment variable not set")
        return {"error": "Missing API key"}
    
    client = anthropic.Anthropic(api_key=api_key)
    
    results = {
        "timestamp": datetime.utcnow().isoformat() + "Z",
        "urls_processed": len(urls),
        "docs_created": [],
        "errors": []
    }
    
    # Process each URL
    for url in urls:
        print(f"📄 Processing: {url}")
        
        prompt = f"""Please fetch and process this documentation URL: {url}

Extract and format the content as clean markdown with:
- Source URL at the top
- Fetch timestamp
- All key sections preserved
- Code examples properly formatted
- Important notes and warnings included

Return ONLY the markdown content, no additional commentary."""

        try:
            message = client.messages.create(
                model="claude-3-5-sonnet-20241022",
                max_tokens=8192,
                temperature=0,
                system="You are a documentation processor. Extract and format documentation content clearly and accurately.",
                messages=[
                    {"role": "user", "content": prompt}
                ]
            )
            
            content = message.content[0].text
            
            # Add header if not present
            if not content.startswith("**Source URL:**"):
                header = f"""**Source URL:** {url}
**Fetch Timestamp:** {datetime.utcnow().isoformat()}Z

"""
                content = header + content
            
            # Save to file
            filename = slugify(url)
            filepath = Path(docs_dir) / filename
            filepath.parent.mkdir(parents=True, exist_ok=True)
            
            with open(filepath, 'w', encoding='utf-8') as f:
                f.write(content)
            
            print(f"   ✅ Saved to: {filepath}")
            results["docs_created"].append(str(filepath))
            
        except Exception as e:
            print(f"   ❌ Error: {str(e)}")
            results["errors"].append({
                "url": url,
                "error": str(e)
            })
    
    return results

def main():
    """Main execution function."""
    
    # Parse arguments
    docs_dir = sys.argv[1] if len(sys.argv) > 1 else "docs"
    urls_file = sys.argv[2] if len(sys.argv) > 2 else "urls.txt"
    
    print(f"📅 Documentation Refresh - {datetime.now()}")
    print("━" * 40)
    
    # Create directories
    Path(docs_dir).mkdir(parents=True, exist_ok=True)
    Path("logs").mkdir(parents=True, exist_ok=True)
    
    # Read URLs
    if not Path(urls_file).exists():
        # Create default urls.txt
        default_urls = """# Documentation URLs to refresh daily
# Add one URL per line
https://docs.anthropic.com/en/api/openai-sdk
https://docs.anthropic.com/en/api/messages
https://docs.anthropic.com/en/api/client-sdks
https://docs.anthropic.com/en/api/streaming
"""
        with open(urls_file, 'w') as f:
            f.write(default_urls)
        print(f"✅ Created {urls_file} with default URLs")
    
    # Read and filter URLs
    with open(urls_file, 'r') as f:
        urls = [
            line.strip() 
            for line in f 
            if line.strip() and not line.strip().startswith('#')
        ]
    
    if not urls:
        print(f"❌ No URLs found in {urls_file}")
        sys.exit(1)
    
    print(f"📋 Found {len(urls)} URLs to refresh")
    print()
    
    # Fetch documentation
    results = fetch_with_anthropic(urls, docs_dir)
    
    # Save results
    timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
    json_output = f"logs/refresh_{timestamp}.json"
    
    with open(json_output, 'w') as f:
        json.dump(results, f, indent=2)
    
    # Summary
    print()
    print("━" * 40)
    print("✅ Refresh complete")
    print(f"📊 Processed: {results['urls_processed']} URLs")
    print(f"✅ Created: {len(results['docs_created'])} documents")
    if results['errors']:
        print(f"❌ Errors: {len(results['errors'])}")
    print(f"📊 Results saved to: {json_output}")
    
    # Show created files
    if results['docs_created']:
        print()
        print("📁 Documents created:")
        for doc in results['docs_created'][:5]:
            print(f"   - {doc}")
        if len(results['docs_created']) > 5:
            print(f"   ... and {len(results['docs_created']) - 5} more")
    
    # Exit with error code if there were failures
    sys.exit(1 if results.get('errors') else 0)

if __name__ == "__main__":
    main()