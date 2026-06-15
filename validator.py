import os
import csv

def validate_gtfs(gtfs_dir):
    print("=========================================")
    print("🚌 STARTING TAKORADI TROTRO DATA VALIDATION")
    print("=========================================\n")
    
    errors = 0
    
    # 1. Check if crucial files exist
    required_files = ['agency.txt', 'routes.txt', 'stops.txt', 'fare_attributes.txt', 'fare_rules.txt']
    for file in required_files:
        path = os.path.join(gtfs_dir, file)
        if not os.path.exists(path):
            print(f"❌ CRITICAL ERROR: Missing required file: {file}")
            errors += 1
            return

    # 2. Parse Route IDs into memory for validation cross-checking
    valid_route_ids = set()
    routes_path = os.path.join(gtfs_dir, 'routes.txt')
    with open(routes_path, mode='r', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for row in reader:
            valid_route_ids.add(row['route_id'])

    # 3. Validate Stops for coordinate validity
    stops_path = os.path.join(gtfs_dir, 'stops.txt')
    print("📋 Checking stops.txt coordinates...")
    with open(stops_path, mode='r', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for line_num, row in enumerate(reader, start=2):
            name = row.get('stop_name', 'Unknown')
            lat = row.get('stop_lat')
            lon = row.get('stop_lon')
            
            if not lat or not lon:
                print(f"  ❌ Line {line_num}: Stop '{name}' is missing geo-coordinates.")
                errors += 1
                continue
                
            try:
                # Double check that coordinates are actual numbers
                float(lat)
                float(lon)
            except ValueError:
                print(f"  ❌ Line {line_num}: Stop '{name}' has corrupt invalid numbers for coordinates.")
                errors += 1

    # 4. Validate Fare Rules reference valid Route IDs
    rules_path = os.path.join(gtfs_dir, 'fare_rules.txt')
    print("\n💵 Checking fare_rules.txt links...")
    with open(rules_path, mode='r', encoding='utf-8') as f:
        reader = csv.DictReader(f)
        for line_num, row in enumerate(reader, start=2):
            r_id = row.get('route_id')
            if r_id and r_id not in valid_route_ids:
                print(f"  ❌ Line {line_num}: Fare rule references a route_id that does not exist: '{r_id}'")
                errors += 1

    # Final Evaluation Summary
    print("\n=========================================")
    if errors == 0:
        print("✅ SUCCESS: All Takoradi transit files look structurally sound!")
    else:
        print(f"❌ FAILED: Found {errors} validation errors. Fix them before publishing.")
    print("=========================================")

if __name__ == "__main__":
    # Point directly to your workspace data directory
    target_directory = os.path.join("data", "gtfs")
    validate_gtfs(target_directory)