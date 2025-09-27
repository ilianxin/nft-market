// SPDX-License-Identifier: MIT
pragma solidity ^0.8.19;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/token/ERC721/extensions/ERC721URIStorage.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
import "@openzeppelin/contracts/utils/Counters.sol";

/**
 * @title 测试NFT合约
 * @dev 用于测试的简单ERC721合约
 */
contract TestNFT is ERC721, ERC721URIStorage, Ownable {
    using Counters for Counters.Counter;
    Counters.Counter private _tokenIdCounter;
    
    // 映射tokenId到创建者
    mapping(uint256 => address) public creators;
    
    event NFTMinted(uint256 indexed tokenId, address indexed creator, address indexed owner, string tokenURI);
    
    constructor() ERC721("Test NFT Collection", "TNC") {}
    
    /**
     * @dev 铸造NFT
     * @param to NFT接收者地址
     * @param tokenURI NFT元数据URI
     * @return tokenId 新铸造的NFT ID
     */
    function mint(address to, string memory tokenURI) external returns (uint256) {
        _tokenIdCounter.increment();
        uint256 tokenId = _tokenIdCounter.current();
        
        creators[tokenId] = msg.sender;
        _mint(to, tokenId);
        _setTokenURI(tokenId, tokenURI);
        
        emit NFTMinted(tokenId, msg.sender, to, tokenURI);
        return tokenId;
    }
    
    /**
     * @dev 批量铸造NFT（用于测试）
     * @param to NFT接收者地址
     * @param count 铸造数量
     * @param baseURI 基础URI
     */
    function mintBatch(address to, uint256 count, string memory baseURI) external returns (uint256[] memory) {
        uint256[] memory tokenIds = new uint256[](count);
        
        for (uint256 i = 0; i < count; i++) {
            _tokenIdCounter.increment();
            uint256 tokenId = _tokenIdCounter.current();
            
            creators[tokenId] = msg.sender;
            _mint(to, tokenId);
            
            string memory tokenURI = string(abi.encodePacked(baseURI, "/", Strings.toString(tokenId), ".json"));
            _setTokenURI(tokenId, tokenURI);
            
            tokenIds[i] = tokenId;
            emit NFTMinted(tokenId, msg.sender, to, tokenURI);
        }
        
        return tokenIds;
    }
    
    /**
     * @dev 获取当前Token ID计数器
     */
    function getCurrentTokenId() external view returns (uint256) {
        return _tokenIdCounter.current();
    }
    
    /**
     * @dev 获取NFT创建者
     */
    function getCreator(uint256 tokenId) external view returns (address) {
        return creators[tokenId];
    }
    
    /**
     * @dev 检查Token是否存在
     */
    function exists(uint256 tokenId) external view returns (bool) {
        return _exists(tokenId);
    }
    
    // 重写必需的函数
    function tokenURI(uint256 tokenId) public view override(ERC721, ERC721URIStorage) returns (string memory) {
        return super.tokenURI(tokenId);
    }
    
    function _burn(uint256 tokenId) internal override(ERC721, ERC721URIStorage) {
        super._burn(tokenId);
    }
    
    function supportsInterface(bytes4 interfaceId) public view override(ERC721, ERC721URIStorage) returns (bool) {
        return super.supportsInterface(interfaceId);
    }
}
