import os
import tempfile
from analyzer.vector_store import search_evidence

def test_search_evidence_empty():
    case_id = "non_existent_case_999"
    result = search_evidence(case_id, "search query")
    assert result["status"] == "success"
    assert result["results"] == []
