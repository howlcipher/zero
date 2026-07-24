import time
import os
import sys
import json
from openai import OpenAI

def tail_file(filepath):
    """Generator that yields new lines from a file as they are written."""
    with open(filepath, 'r') as f:
        f.seek(0, 2)  # Go to the end of the file
        while True:
            line = f.readline()
            if not line:
                time.sleep(0.1)
                continue
            yield line

def check_anomalies(client, events):
    prompt = "You are an AI observability agent monitoring telemetry from a Go application. Review the following execution trace and flag any anomalous behavior, infinite loops, or unexpected state changes. If everything looks normal, respond with 'NORMAL'. Otherwise, explain the anomaly.\n\nEvents:\n"
    for e in events:
        prompt += e + "\n"
        
    try:
        response = client.chat.completions.create(
            model="llama3",
            messages=[{"role": "user", "content": prompt}],
            max_tokens=200
        )
        reply = response.choices[0].message.content.strip()
        if "NORMAL" not in reply.upper() or len(reply) > 20:
            print(f"\n[⚠️ ANOMALY DETECTED]\n{reply}\n")
        else:
            print("[OK] Trace segment looks normal.")
    except Exception as e:
        print(f"[Error] LLM request failed: {e}")

def main():
    telemetry_file = "telemetry.jsonl"
    print(f"Starting Standalone Observer Agent. Tailing {telemetry_file}...")
    
    if not os.path.exists(telemetry_file):
        open(telemetry_file, 'a').close()
        
    client = OpenAI(base_url="http://localhost:11434/v1", api_key="ollama")
    
    buffer = []
    max_buffer_size = 10
    
    try:
        for line in tail_file(telemetry_file):
            line = line.strip()
            if line:
                print(f"Observed: {line}")
                buffer.append(line)
                
            if len(buffer) >= max_buffer_size:
                print("Analyzing recent trace segment...")
                check_anomalies(client, buffer)
                buffer = []
    except KeyboardInterrupt:
        print("\nObserver agent stopped.")
        sys.exit(0)

if __name__ == "__main__":
    main()
