import chromadb
from chromadb.utils import embedding_functions
import os
import json

# Local-first embedding function definition:
# We instantiate SentenceTransformerEmbeddingFunction which runs a local transformer model.
# By default, we use 'all-MiniLM-L6-v2', a fast, high-performance, and lightweight model (under 120MB)
# that converts sentences/paragraphs into 384-dimensional dense vectors.
# This runs completely offline on the user's CPU/GPU and does not require any external internet requests.
default_ef = embedding_functions.SentenceTransformerEmbeddingFunction(model_name="all-MiniLM-L6-v2")

def get_client(case_id):
    """
    Initializes and returns a persistent ChromaDB database client for a specific case.
    
    Parameters:
    - case_id (str): The unique ID of the investigation case.
    
    The database state is saved to a hidden folder '.vector_store' located inside the case's 
    evidence directory (e.g. evidence_storage/<case_id>/.vector_store). This keeps case data 
    isolated and local.
    """
    persist_directory = os.path.join("evidence_storage", case_id, ".vector_store")
    os.makedirs(persist_directory, exist_ok=True)
    # Instantiate a persistent Chroma client. It writes database files to disk at the path.
    return chromadb.PersistentClient(path=persist_directory)

def index_evidence(case_id, evidence_files):
    """
    Indexes a list of evidence files into the case's vector store database for semantic retrieval.
    
    Parameters:
    - case_id (str): The unique ID of the case.
    - evidence_files (list of str): Paths to raw evidence files (e.g., ports.json, whois.json).
    
    Returns:
    - dict: A success status dictionary indicating the count of indexed files.
    """
    # 1. Resolve case-specific database client.
    client = get_client(case_id)
    
    # 2. Get or create the vector collection named 'evidence'.
    # We pass default_ef to ensure Chroma calls sentence-transformers to calculate embedding vectors
    # when documents are added or queried.
    collection = client.get_or_create_collection(
        name="evidence", 
        embedding_function=default_ef
    )

    ids = []
    documents = []
    metadatas = []

    # 3. Iterate over each provided file path and extract contents.
    for file_path in evidence_files:
        # Check if the file exists to prevent throwing runtime exceptions.
        if not os.path.exists(file_path):
            continue
            
        file_name = os.path.basename(file_path)
        
        # Operational security & layout check:
        # Skip hidden files/directories (starting with '.') or HTML visual reports.
        if file_name.startswith(".") or file_name == "report.html":
            continue

        # Read the raw file content. We use 'ignore' errors flag to handle binary files or
        # non-UTF-8 characters gracefully without crashing the analysis run.
        with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
            content = f.read()
            
        # Store file data in memory for bulk upsert.
        # Currently, the entire file content is indexed as a single document.
        # Document ID is set to the file name.
        ids.append(file_name)
        documents.append(content)
        metadatas.append({"path": file_path, "name": file_name})

    # 4. If we gathered any valid files, upsert them to the database.
    # collection.upsert inserts the document if it's new, or updates it if the file ID already exists.
    if ids:
        collection.upsert(
            ids=ids,
            documents=documents,
            metadatas=metadatas
        )
        return {"status": "success", "indexed_count": len(ids)}
    
    return {"status": "success", "indexed_count": 0}

def search_evidence(case_id, query, n_results=3):
    """
    Performs a semantic vector search across the indexed evidence documents.
    
    Parameters:
    - case_id (str): The unique ID of the case to search.
    - query (str): The natural language query (e.g. "email addresses found").
    - n_results (int): Max number of matching evidence files to retrieve.
    
    Returns:
    - dict: A status dictionary containing formatted matches, each with the file name, content,
            metadata, and vector distance (similarity score).
    """
    # 1. Resolve case-specific database client.
    client = get_client(case_id)
    
    # 2. Retrieve the collection. It throws an exception if the collection has not been created yet.
    collection = client.get_collection(name="evidence", embedding_function=default_ef)
    
    # 3. Query the collection. The text query is converted into a vector embedding under the hood,
    # and Chroma performs cosine similarity matching to find the closest vectors.
    results = collection.query(
        query_texts=[query],
        n_results=n_results
    )
    
    # 4. Format search results into a clean list of dictionaries for Go to parse.
    formatted_results = []
    # Chroma returns lists of results. results['ids'][0] holds matching IDs for the first query text.
    for i in range(len(results['ids'][0])):
        formatted_results.append({
            "id": results['ids'][0][i],
            "content": results['documents'][0][i],
            "metadata": results['metadatas'][0][i],
            "distance": results['distances'][0][i]
        })
        
    return {"status": "success", "results": formatted_results}
