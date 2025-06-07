import json
import numpy as np  # type: ignore 
from sentence_transformers import SentenceTransformer  # type: ignore
from sklearn.cluster import KMeans  # type: ignore
from sklearn.feature_extraction.text import TfidfVectorizer  # type: ignore 
import string  
import re  

def cluster_words(input_file, output_file, n_clusters=10, n_difficulties=3):  
    """
    Cluster words based on their semantic meaning and assign difficulty levels.
    
    Args:
        input_file: Path to input file containing one word per line
        output_file: Path to output JSON file
        n_clusters: Number of semantic clusters to create
        n_difficulties: Number of difficulty levels to assign (now 3: easy, medium, hard)  
    """
    # Step 1: Load words from file
    with open(input_file, 'r') as f:
        words = [line.strip() for line in f if line.strip()]

    # Step 2: Load SentenceTransformer model
    model = SentenceTransformer('all-MiniLM-L6-v2')
    embeddings = model.encode(words)

    # Step 3: Perform KMeans clustering
    kmeans = KMeans(n_clusters=n_clusters, random_state=42)
    kmeans.fit(embeddings)
    
    # Step 4: Calculate difficulty scores based on multiple factors
    difficulty_scores = calculate_difficulty_scores(words)
    
    # Step 5: Assign difficulty levels (1-3) based on percentiles  
    difficulty_bins = np.linspace(0, 100, n_difficulties+1)
    difficulty_percentiles = np.percentile(difficulty_scores, difficulty_bins)
    difficulty_levels = np.digitize(difficulty_scores, difficulty_percentiles[1:-1]) + 1  # +1 to start from level 1
    
    # Step 6: Create the new JSON structure with difficulty labels  
    result = []
    difficulty_labels = ["easy", "medium", "hard"]  # added difficulty labels
    for i, word in enumerate(words):
        cluster_label = int(kmeans.labels_[i])
        difficulty_level = int(difficulty_levels[i])
        difficulty_label = difficulty_labels[difficulty_level-1]  # get the label based on level
        result.append({
            "word": word,
            "category": f"Cluster {cluster_label}",
            "difficulty": difficulty_level,
            "difficulty_label": difficulty_label  
        })

    # Step 7: Save to JSON file in the requested format
    with open(output_file, 'w') as f:
        json.dump(result, f, indent=2)

    print(f"✅ Clustered {len(words)} into {n_clusters} clusters with {n_difficulties} difficulty levels (easy, medium, hard). Output saved to {output_file}")  # updated success message

# New helper function to calculate word difficulty
def calculate_difficulty_scores(words):
    """
    Calculate difficulty scores for words based on multiple factors:
    1. Word length
    2. Character complexity (presence of unusual characters)
    3. Word rarity (inversely related to frequency)
    4. Syllable count
    
    Returns an array of difficulty scores
    """
    scores = []
    
    # Get word length scores
    length_scores = np.array([len(word) for word in words])
    
    # Character complexity (number of non-alphanumeric characters)
    complexity_scores = np.array([sum(1 for char in word if char not in string.ascii_letters) for word in words])
    
    # Word rarity using TF-IDF (higher values for rare words)
    tfidf = TfidfVectorizer(analyzer='char', ngram_range=(2, 3))
    rarity_features = tfidf.fit_transform(words).toarray()
    rarity_scores = np.mean(rarity_features, axis=1)
    
    # Syllable count estimation (rough heuristic)
    def count_syllables(word):
        word = word.lower()
        # Remove trailing e's as they're often silent
        if word.endswith('e'):
            word = word[:-1]
        # Count vowel groups
        count = len(re.findall(r'[aeiouy]+', word))
        return max(1, count)  # Ensure at least one syllable
    
    syllable_scores = np.array([count_syllables(word) for word in words])
    
    # Normalize each score component (min-max scaling)
    def normalize(arr):
        min_val = np.min(arr)
        max_val = np.max(arr)
        if max_val > min_val:
            return (arr - min_val) / (max_val - min_val)
        return np.zeros_like(arr)
    
    norm_length = normalize(length_scores)
    norm_complexity = normalize(complexity_scores)
    norm_rarity = normalize(rarity_scores)
    norm_syllables = normalize(syllable_scores)
    
    # Combine scores with weights (can be adjusted)
    weights = [0.25, 0.25, 0.3, 0.2]  # Length, complexity, rarity, syllables
    combined_scores = (
        weights[0] * norm_length +
        weights[1] * norm_complexity +
        weights[2] * norm_rarity +
        weights[3] * norm_syllables
    )
    
    return combined_scores

# Example usage
cluster_words("oxford3000_clean.txt", "clustered_with_difficulty.json", n_clusters=10, n_difficulties=3)  