// Pacote
package main

// Importações
import (
	"net/http"
	"uuid"

	"github.com/gin-gonic/gin"
)

// Struct Produto
type Produto struct {
	Codigo    string  `json:"codigo"`
	Nome      string  `json:"nome"`
	Categoria string  `json:"categoria"`
	Valor     float64 `json:"valor"`
}

// Lista (Slice)
var produtos = []Produto{}

// Função principal
func main() {

	// Roteador
	r := gin.Default()

	// Configurar CORS
	r.Use(corsMiddleware())

	// Rotas
	r.GET("/", helloWorld)
	r.GET("/produtos", listarProdutos)
	r.POST("/produtos", cadastrarProduto)
	r.PUT("/produtos/:codigo", alterarProduto)
	r.DELETE("/produtos/:codigo", removerProduto)

	// Configurar o servidor
	r.Run(":8080")

}

// Middleware de CORS adaptado para o Gin
func corsMiddleware() gin.HandlerFunc {
	// Retorno
	return func(c *gin.Context) {

		// Configurações
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// Options
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		// Next
		c.Next()
	}
}

// Função para retornar um Hello World
func helloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"mensagem": "Hello World!"})
}

// Função para listar todos os produtos
func listarProdutos(c *gin.Context) {
	c.JSON(http.StatusOK, produtos)
}

// Função para cadastrar um novo produto
func cadastrarProduto(c *gin.Context) {
	// Referência do struct Produto
	var novoProduto Produto

	// Extrair informações do JSON e atribuir na referência novoProduto
	if err := c.ShouldBindJSON(&novoProduto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	// Gerar o UUID
	novoProduto.Codigo = uuid.New().String()

	// Adicionar o novoProduto a lista
	produtos = append(produtos, novoProduto)

	// Retorno do endpoint ao cadastrar um novo produto
	c.JSON(http.StatusCreated, novoProduto)
}

// Função para atualizar um produto específico
func alterarProduto(c *gin.Context) {
	// Extrair o código do produto via URL
	codigo := c.Param("codigo")

	// Referência do struct Produto
	var produtoAtualizado Produto

	// Extrair informações do JSON e atribuir na referência produtoAtualizado
	if err := c.ShouldBindJSON(&produtoAtualizado); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
		return
	}

	// Laço de repetição para localizar o produto na lista e modificar
	for indice, produto := range produtos {
		// Condicional
		if produto.Codigo == codigo {
			produtoAtualizado.Codigo = codigo
			produtos[indice] = produtoAtualizado

			c.JSON(http.StatusOK, produtoAtualizado)
			return
		}
	}

	// Retorno do endpoint caso o produto não seja encontrado
	c.JSON(http.StatusNotFound, gin.H{"mensagem": "Produto não encontrado."})
}

// Função para remover um produto específico
func removerProduto(c *gin.Context) {
	// Extrair o código do produto via URL
	codigo := c.Param("codigo")

	// Laço de repetição para localizar o produto na lista e remover
	for indice, produto := range produtos {
		// Condicional
		if produto.Codigo == codigo {

			produtos = append(produtos[:indice], produtos[indice+1:]...)

			c.JSON(http.StatusOK, gin.H{"mensagem": "Produto removido."})
			return
		}
	}

	// Retorno do endpoint caso o produto não seja encontrado
	c.JSON(http.StatusNotFound, gin.H{"mensagem": "Produto não encontrado."})
}
