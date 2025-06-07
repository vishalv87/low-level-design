import json

def merge_json_data(words_file, clusters_file, difficulty_levels_file, output_file):
    """
    Merges data from three JSON files into a unified structure.
    
    Args:
        words_file: Path to the file containing words with their category and difficulty
        clusters_file: Path to the file containing cluster metadata
        difficulty_levels_file: Path to the file containing difficulty level metadata
        output_file: Path to save the combined output
    """
    # Load data from each file
    with open(words_file, 'r', encoding='utf-8') as f:
        words_data = json.load(f)
    
    with open(clusters_file, 'r', encoding='utf-8') as f:
        clusters_data = json.load(f)
    
    with open(difficulty_levels_file, 'r', encoding='utf-8') as f:
        difficulty_data = json.load(f)
    
    # Create the combined structure
    result = {
        "words": words_data,  # This is already an array of word objects
        "clusters": clusters_data["clusters"],  # Extract just the clusters array
        "difficulty_levels": difficulty_data["difficulty_levels"]  # Extract just the difficulty_levels array
    }
    
    # For each word, update the category field to use the primary name instead of "Cluster X"
    cluster_map = {f"Cluster {c['cluster_id']}": c["primary_name"] for c in clusters_data["clusters"]}
    
    for word in result["words"]:
        if word["category"] in cluster_map:
            word["category"] = cluster_map[word["category"]]
    
    # Save the combined JSON to the output file
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(result, f, indent=2)
    
    print(f"✅ Successfully merged data into {output_file}")


# Alternative version that can work with just the cluster number
def merge_json_data_alt(words_file, clusters_file, difficulty_levels_file, output_file):
    """
    Alternative version that works if the category is stored as just a number or "Cluster X"
    """
    # Load data from each file
    with open(words_file, 'r', encoding='utf-8') as f:
        words_data = json.load(f)
    
    with open(clusters_file, 'r', encoding='utf-8') as f:
        clusters_data = json.load(f)
    
    with open(difficulty_levels_file, 'r', encoding='utf-8') as f:
        difficulty_data = json.load(f)
    
    # Create a mapping from cluster IDs to their primary names
    cluster_map = {}
    for cluster in clusters_data["clusters"]:
        cluster_id = cluster["cluster_id"]
        # Handle both "Cluster X" and just the number X
        cluster_map[f"Cluster {cluster_id}"] = cluster["primary_name"]
        cluster_map[str(cluster_id)] = cluster["primary_name"]
        cluster_map[cluster_id] = cluster["primary_name"]
    
    # For each word, update the category using the primary name
    for word in words_data:
        if word["category"] in cluster_map:
            word["category"] = cluster_map[word["category"]]
    
    # Create the combined structure
    result = {
        "words": words_data,
        "clusters": clusters_data["clusters"],
        "difficulty_levels": difficulty_data["difficulty_levels"]
    }
    
    # Save the combined JSON to the output file
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(result, f, indent=2)
    
    print(f"✅ Successfully merged data into {output_file}")


# Add a function that also adds the difficulty label to each word
def merge_json_data_with_labels(words_file, clusters_file, difficulty_levels_file, output_file):
    """
    Merges data and adds difficulty labels to each word based on the numeric difficulty level.
    """
    # Load data from each file
    with open(words_file, 'r', encoding='utf-8') as f:
        words_data = json.load(f)
    
    with open(clusters_file, 'r', encoding='utf-8') as f:
        clusters_data = json.load(f)
    
    with open(difficulty_levels_file, 'r', encoding='utf-8') as f:
        difficulty_data = json.load(f)
    
    # Create mappings
    cluster_map = {f"Cluster {c['cluster_id']}": c["primary_name"] for c in clusters_data["clusters"]}
    difficulty_map = {d["level"]: d["label"] for d in difficulty_data["difficulty_levels"]}
    
    # Update each word with the proper category name and add difficulty label
    for word in words_data:
        # Update category if it matches any key in the cluster map
        if word["category"] in cluster_map:
            word["category"] = cluster_map[word["category"]]
        
        # Add difficulty label based on the numeric difficulty level
        difficulty_level = word["difficulty"]
        if isinstance(difficulty_level, (int, str)) and int(difficulty_level) in difficulty_map:
            word["difficulty_label"] = difficulty_map[int(difficulty_level)]
    
    # Create the combined structure
    result = {
        "words": words_data,
        "clusters": clusters_data["clusters"],
        "difficulty_levels": difficulty_data["difficulty_levels"]
    }
    
    # Save the combined JSON to the output file
    with open(output_file, 'w', encoding='utf-8') as f:
        json.dump(result, f, indent=2)
    
    print(f"✅ Successfully merged data into {output_file}")


# Example usage
if __name__ == "__main__":
    # Define your file paths here
    words_file = "clustered_with_difficulty.json"
    clusters_file = "cluster_data.json"
    difficulty_file = "difficulty_levels.json"
    output_file = "combined_vocabulary.json"
    
    # Choose which function to use based on your needs
    merge_json_data_with_labels(words_file, clusters_file, difficulty_file, output_file)
    
    # Or, if you want to include difficulty labels:
    # merge_json_data_with_labels(words_file, clusters_file, difficulty_file, output_file)