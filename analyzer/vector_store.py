import chromadb
from chromadb.utils import embedding_functions
import os
import json

# Local-first embedding function
# This uses the sentence-transformers library internally
default_ef = embedding_functions.SentenceTransformerEmbeddingFunction(model_name="all-MiniLM-L6-v2")

def get_client(case_id):
    """
    Initializes a persistent ChromaDB client for a specific case.
    """
    persist_directory = os.path.join("evidence_storage", case_id, ".vector_store")
    os.makedirs(persist_directory, exist_ok=True)
    return chromadb.PersistentClient(path=persist_directory)

def index_evidence(case_id, evidence_files):
    """
    Indexes a list of evidence files into the case's vector store.
    """
    client = get_client(case_id)
    collection = client.get_or_create_collection(
        name="evidence", 
        embedding_function=default_ef
    )

    ids = []
    documents = []
    metadatas = []

    for file_path in evidence_files:
        if not os.path.exists(file_path):
            continue
            
        file_name = os.path.basename(file_path)
        
        # Skip internal files
        if file_name.startswith(".") or file_name == "report.html":
            continue

        with open(file_path, 'r', encoding='utf-8', errors='ignore') as f:
            content = f.read()
            
        # For JSON files, we might want to prettify or extract specific fields
        # but for now, we index the whole content.
        
        ids.append(file_name)
        documents.append(content)
        metadatas.append({"path": file_path, "name": file_name})

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
    Performs a semantic search across the indexed evidence.
    """
    client = get_client(case_id)
    collection = client.get_collection(name="evidence", embedding_function=default_ef)
    
    results = collection.query(
        query_texts=[query],
        n_results=n_results
    )
    
    formatted_results = []
    for i in range(len(results['ids'][0])):
        formatted_results.append({
            "id": results['ids'][0][i],
            "content": results['documents'][0][i],
            "metadata": results['metadatas'][0][i],
            "distance": results['distances'][0][i]
        })
        
    return {"status": "success", "results": formatted_results}
