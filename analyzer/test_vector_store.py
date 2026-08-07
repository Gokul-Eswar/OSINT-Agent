import os
import tempfile
from analyzer.vector_store import search_evidence, index_evidence

def test_search_evidence_empty():
    case_id = "non_existent_case_999"
    result = search_evidence(case_id, "search query")
    assert result["status"] == "success"
    assert result["results"] == []

def test_index_evidence_nonexistent_files():
    case_id = "test_case_empty_files"
    result = index_evidence(case_id, ["/non_existent_path/foo.json"])
    assert result["status"] == "success"
    assert result["indexed_count"] == 0
